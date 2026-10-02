#!/usr/bin/env bash
# Live lifecycle test of terraform-provider-elestio against the real Elestio API.
#
# WHAT IT DOES
#   Builds the provider from this checkout, then drives the real Terraform CLI to create
#   ONE small billed service, update it, delete it outside Terraform, destroy it, and
#   create/destroy it once more. Every step prints PASS / FAIL / INCONCLUSIVE.
#
# COST AND SAFETY
#   hetzner / nbg1 / MEDIUM-2C-4G is about $0.03 per hour and the service lives for a few
#   minutes. The script destroys what it creates. If it is interrupted, run:
#       scripts/live-lifecycle-test.sh cleanup
#   The token and JWT are never printed. Everything written to the output folder is
#   redacted, and the script verifies no secret is left in it.
#
# SIGN-IN BUDGET (handled automatically)
#   The provider now reuses one JWT for the whole run (and caches it under $WORK/cache), so a
#   full run needs only about 6 sign-ins: the pre-flight, the provider's first sign-in, and
#   one for each API helper call. The script still keeps a ledger ($OUT/signin-ledger.txt) and
#   sleeps before any helper call that would exceed SIGNIN_CAP (default 12 of 15 per hour).
#   The ledger cannot see sign-ins made by other tools; do not run other Elestio commands
#   while it runs.
# SIGN-IN BUDGET (background)
#   The backend allows 15 sign-ins per user per hour (and per IP), and every terraform
#   command that configures the provider signs in once. A full run now uses about 6. Start it
#   when no other Elestio API or Terraform commands ran in the last hour. Use
#   WAIT_FOR_WINDOW=1 to make the script retry the pre-flight sign-in every 10 minutes
#   (up to 90 minutes) instead of stopping.
#
# USAGE
#   export ELESTIO_EMAIL=you@example.com ELESTIO_API_TOKEN=... PROJECT_ID=<project id>
#   scripts/live-lifecycle-test.sh              # full test
#   scripts/live-lifecycle-test.sh preflight    # sign-in check only (1 sign-in)
#   scripts/live-lifecycle-test.sh cleanup      # destroy leftovers from an earlier run
#
# OPTIONS (environment)
#   TEST_FIREWALL=1  also test firewall rules on the same service (add, replace, disable,
#                    re-enable, and invalid/incomplete rules rejected). Adds about 4 minutes.
#   PROJECT_ID   project to use (REQUIRED, the project must exist in your account)
#   PROVIDER_NAME / DATACENTER / SERVER_TYPE   (default hetzner / nbg1 / MEDIUM-2C-4G)
#   OUT_DIR      output folder                 (default ./live-test-output)
#   WAIT_FOR_WINDOW=1  retry the pre-flight while the rate limit is active
#
# SEND ME: OUT_DIR/report.txt  and  OUT_DIR/terraform.redacted.log
#   (preflight and cleanup runs write report-<mode>.txt and terraform-<mode>.redacted.log instead)
set -u

REPO="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MODE="${1:-full}"
PROJECT="${PROJECT_ID:?export PROJECT_ID=<an Elestio project id you can create services in>}"
PROVIDER_NAME="${PROVIDER_NAME:-hetzner}"
DATACENTER="${DATACENTER:-nbg1}"
SERVER_TYPE="${SERVER_TYPE:-MEDIUM-2C-4G}"
OUT="${OUT_DIR:-$PWD/live-test-output}"
WORK="$OUT/work"
RAW="$OUT/.raw"            # unredacted logs; deleted at the end
LOG="$RAW/terraform.log"
DEBUG="$RAW/provider.debug"
# preflight and cleanup runs write their own files so they never overwrite the evidence of a
# full run.
SUFFIX=""; [ "$MODE" != "full" ] && SUFFIX="-$MODE"
REPORT="$OUT/report$SUFFIX.txt"
PW="AuditPass12345"

: "${ELESTIO_EMAIL:?export ELESTIO_EMAIL first}"
: "${ELESTIO_API_TOKEN:?export ELESTIO_API_TOKEN first}"
for bin in terraform go python3; do
  command -v "$bin" >/dev/null || { echo "missing required tool: $bin"; exit 1; }
done

mkdir -p "$WORK" "$RAW" && chmod 700 "$RAW"
warn_leftover() {
  if [ -d "$WORK/.terraform" ] && (cd "$WORK" && terraform state list 2>/dev/null | grep -q elestio); then
    echo "!! A billed service may still exist. Run: OUT_DIR=$OUT $0 cleanup   (or delete it in the dashboard)"
  fi
}
trap warn_leftover EXIT

# ---------------------------------------------------------------- API helper
HELPER="$RAW/liveapi.py"
cat > "$HELPER" <<'PY'
import json, os, re, sys, time, urllib.request, urllib.error
BASE = "https://api.elest.io"

def post(path, body):
    req = urllib.request.Request(BASE + path, data=json.dumps(body).encode(),
                                 headers={"Content-Type": "application/json"})
    try:
        r = urllib.request.urlopen(req, timeout=120)
        return r.status, r.read().decode()
    except urllib.error.HTTPError as e:
        return e.code, e.read().decode()

def sign_in():
    s, b = post("/api/auth/checkAPIToken",
                {"email": os.environ["ELESTIO_EMAIL"], "token": os.environ["ELESTIO_API_TOKEN"]})
    try: d = json.loads(b)
    except ValueError: d = {}
    if s != 200 or not d.get("jwt"):
        print("AUTH-FAIL status=%s code=%s" % (s, d.get("code"))); sys.exit(2)
    return d["jwt"]

cmd = sys.argv[1]
jwt = sign_in()
if cmd == "preflight":
    print("AUTH-OK")
elif cmd == "count":
    # Count ACTIVE services whose record mentions the name (inactive 'deleting' records
    # of earlier test services are expected and are not leftovers).
    s, b = post("/api/servers/getServices",
                {"projectId": sys.argv[2], "appid": "", "isActiveService": True, "jwt": jwt})
    if s != 200: print("LIST-FAIL status=%s" % s); sys.exit(2)
    rows = json.loads(b).get("servers", [])
    print("COUNT %d" % sum(1 for r in rows if sys.argv[3] in json.dumps(r)))
elif cmd == "delete_and_wait":
    project, vmid = sys.argv[2], sys.argv[3]
    s, b = post("/api/servers/deleteServer",
                {"projectID": project, "vmID": vmid, "isDeleteServiceWithBackup": True, "jwt": jwt})
    print("DELETE status=%s" % s)
    if s != 200:
        print("DELETE-BODY " + b.replace(jwt, "<jwt>")[:300]); sys.exit(1)
    gone = re.compile(r'"code"\s*:\s*"(InvalidServer|service_deleted)"')
    for _ in range(60):
        time.sleep(10)
        s, b = post("/api/servers/getServerDetails", {"projectID": project, "vmID": vmid, "jwt": jwt})
        if s == 401 and gone.search(b):
            print("GONE yes code=%s" % gone.search(b).group(1)); sys.exit(0)
    print("GONE no last_status=%s" % s); sys.exit(1)
PY
# --------------------------------------------------------------- sign-in ledger
LEDGER="$OUT/signin-ledger.txt"; touch "$LEDGER"
CAP="${SIGNIN_CAP:-12}"
spend() { # spend <n>: wait until n more sign-ins fit in the rolling hour, then record them
  local n="$1" now used oldest wait
  while :; do
    now=$(date +%s)
    awk -v now="$now" '$1 > now-3600' "$LEDGER" > "$LEDGER.tmp" && mv "$LEDGER.tmp" "$LEDGER"
    used=$(awk '{s+=$2} END{print s+0}' "$LEDGER")
    [ $((used + n)) -le "$CAP" ] && break
    oldest=$(sort -n "$LEDGER" | head -1 | awk '{print $1}')
    wait=$(( oldest + 3600 - now + 5 )); [ "$wait" -lt 5 ] && wait=5
    echo "   [budget] $used/$CAP sign-ins used in the last hour; waiting ${wait}s before the next command"
    sleep "$wait"
  done
  echo "$(date +%s) $n" >> "$LEDGER"
}
liveapi() { spend 1; python3 "$HELPER" "$@"; }

# ------------------------------------------------------------------ reporting
pass=0; fail=0; inconc=0
say()   { echo "$*" | tee -a "$REPORT"; }
ok()    { say "PASS          $*"; pass=$((pass+1)); }
bad()   { say "FAIL          $*"; fail=$((fail+1)); }
maybe() { say "INCONCLUSIVE  $*"; inconc=$((inconc+1)); }
auth_failed() { tail -n 60 "$LOG" | grep -q "failed to sign in"; }
result() { # result <exit-code> <label>
  if [ "$1" = "0" ]; then ok "$2"
  elif auth_failed; then maybe "$2 (sign-in rejected, rate limit?)"
  else bad "$2"; fi
}

export TF_CLI_CONFIG_FILE="$WORK/tfrc" CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1
export TF_LOG=DEBUG TF_LOG_PATH="$DEBUG"
# Keep the provider's cached JWT inside the work folder so the test is isolated and the
# token file can be deleted at the end (Linux; macOS/Windows use their own cache dir).
# The provider cache is opt-in; the test enables it, isolated in this folder, so the whole run
# shares one sign-in.
export XDG_CACHE_HOME="$WORK/cache" ELESTIO_JWT_CACHE=on
tf() {
  # The provider signs in once, caches the JWT (isolated under $WORK/cache) and reuses it for
  # every later command, so terraform commands no longer spend sign-ins. Only the python
  # helper below signs in on its own (1 each).
  terraform "$@" -no-color >> "$LOG" 2>&1
}

finish() {
  # Redact everything into files that are safe to send, then delete the raw logs.
  for pair in "$LOG:terraform$SUFFIX.redacted.log" "$DEBUG:provider$SUFFIX.redacted.log"; do
    src="${pair%%:*}"; dst="$OUT/${pair##*:}"
    [ -f "$src" ] || continue
    python3 - "$src" "$dst" "$PW" <<'PYR'
import os, re, sys
src, dst, pw = sys.argv[1:4]
t = open(src, errors="replace").read()
for secret, repl in ((os.environ["ELESTIO_API_TOKEN"], "<token>"), (pw, "<password>")):
    t = t.replace(secret, repl)
t = re.sub(r"eyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_.-]+", "<jwt>", t)
t = re.sub(r"(jwt=)[^&\" \n]+", r"\1<jwt>", t)
open(dst, "w").write(t)
PYR
  done
  leaked=0
  for f in "$OUT"/*.log "$REPORT" "$WORK"/*.txt; do
    [ -f "$f" ] || continue
    grep -q -F -e "$ELESTIO_API_TOKEN" -e "$PW" "$f" 2>/dev/null && { echo "SECRET IN $f"; leaked=1; }
    grep -Eq 'eyJ[A-Za-z0-9_-]{15,}\.' "$f" 2>/dev/null && { echo "JWT IN $f"; leaked=1; }
  done
  rm -rf "$RAW" "$WORK/cache"
  rm -f "$WORK"/show.txt "$WORK"/plan*.txt
  [ "$leaked" = 0 ] && say "Output folder checked: no token, JWT or password in the files to send." \
                    || say "WARNING: a secret pattern was found in the output. Do not send it."
}

# --------------------------------------------------------------- build provider
mkdir -p "$WORK/bin"
( cd "$REPO" && go build -o "$WORK/bin/terraform-provider-elestio" . ) || { echo "provider build failed"; exit 1; }
cat > "$WORK/tfrc" <<EOF
provider_installation {
  dev_overrides {
    "elestio/elestio" = "$WORK/bin"
  }
  direct {}
}
EOF

# ------------------------------------------------------------------- pre-flight
preflight() {
  local tries=1; [ "${WAIT_FOR_WINDOW:-0}" = 1 ] && tries=9
  for i in $(seq 1 "$tries"); do
    out=$(liveapi preflight); rc=$?
    say "pre-flight sign-in (try $i): $out"
    [ $rc -eq 0 ] && return 0
    [ "$i" -lt "$tries" ] && { say "  rate limited, retrying in 10 minutes..."; sleep 600; }
  done
  return 1
}

: > "$REPORT"
say "Live lifecycle test  $(date -u +%FT%TZ)"
say "provider commit: $(cd "$REPO" && git rev-parse --short HEAD 2>/dev/null)$(cd "$REPO" && git diff --quiet 2>/dev/null || echo '+uncommitted')"
say "terraform: $(terraform version | head -1)   project: $PROJECT   target: $PROVIDER_NAME/$DATACENTER/$SERVER_TYPE"
say ""

if [ "$MODE" = "preflight" ]; then preflight; rc=$?; finish; exit $rc; fi

if [ "$MODE" = "cleanup" ]; then
  cd "$WORK" || exit 1
  if ! terraform state list 2>/dev/null | grep -q elestio; then say "Nothing to clean up in $WORK."; finish; exit 0; fi
  preflight || { say "Cannot sign in now; try again later."; finish; exit 2; }
  tf destroy -auto-approve; result $? "cleanup destroy"
  terraform state list 2>/dev/null | grep -q elestio && bad "state still holds resources, check the dashboard" || ok "state is empty"
  finish; [ "$fail" = 0 ]; exit $?
fi

cd "$WORK" || exit 1
if terraform state list 2>/dev/null | grep -q elestio; then
  echo "A previous run left a service in $WORK. Run '$0 cleanup' first."; exit 1
fi
rm -f main.tf terraform.tfstate*; rm -rf .terraform

NAME="tf-audit-$(date +%H%M%S)"
FW_VARS=""; FW_ATTRS=""
if [ "${TEST_FIREWALL:-0}" = 1 ]; then
FW_VARS=$(cat <<'HCLV'

variable "fw_enabled" {
  type    = bool
  default = true
}

variable "fw_drop_ssh" {
  type    = bool
  default = false
}

variable "fw_extra" {
  type = list(object({ port = string, protocol = string }))
  default = []
}

locals {
  fw_all = [
    { type = "input", port = "22", protocol = "tcp", targets = ["0.0.0.0/0", "::/0"] },
    { type = "input", port = "4242", protocol = "udp", targets = ["0.0.0.0/0", "::/0"] },
    { type = "input", port = "80", protocol = "tcp", targets = ["0.0.0.0/0", "::/0"] },
    { type = "input", port = "443", protocol = "tcp", targets = ["0.0.0.0/0", "::/0"] },
  ]
  fw_base  = [for r in local.fw_all : r if !(var.fw_drop_ssh && r.port == "22")]
  fw_extra = [for r in var.fw_extra : { type = "input", port = r.port, protocol = r.protocol, targets = ["0.0.0.0/0", "::/0"] }]
}
HCLV
)
FW_ATTRS=$(cat <<'HCLA'
  firewall_enabled    = var.fw_enabled
  firewall_user_rules = [for r in concat(local.fw_base, local.fw_extra) : r if var.fw_enabled]
HCLA
)
fi
cat > main.tf <<EOF
terraform {
  required_providers {
    elestio = {
      source = "elestio/elestio"
    }
  }
}

provider "elestio" {
}

variable "alerts" {
  type    = bool
  default = true
}

resource "elestio_ubuntu" "t" {
  project_id       = "$PROJECT"
  server_name      = "$NAME"
  provider_name    = "$PROVIDER_NAME"
  datacenter       = "$DATACENTER"
  server_type      = "$SERVER_TYPE"
  default_password = "$PW"
  alerts_enabled   = var.alerts
$FW_ATTRS
}
$FW_VARS

output "service_id" {
  value = elestio_ubuntu.t.id
}
EOF

say "== 0. pre-flight"
preflight || { say "Cannot sign in now (rate limit?). Nothing was created. Retry later or set WAIT_FOR_WINDOW=1."; finish; exit 2; }

say "== 1. validate"
tf validate && ok "terraform validate" || bad "terraform validate"

say "== 2. create (waits for deployment, up to 20 min)"
tf apply -auto-approve; rc=$?
result $rc "apply (create) succeeded"
if [ $rc -ne 0 ]; then say "Create failed. If a service was created, run: $0 cleanup"; finish; exit 1; fi
SVC_ID=$(terraform output -raw service_id 2>/dev/null)
[ -n "$SVC_ID" ] && ok "service id recorded in state ($SVC_ID)" || bad "service id missing from state"

say "== 3. no drift right after create"
terraform plan -detailed-exitcode -no-color > plan1.txt 2>&1; rc=$?; cat plan1.txt >> "$LOG"
if [ $rc -eq 0 ]; then ok "plan is empty"
elif grep -q "failed to sign in" plan1.txt; then maybe "drift check (sign-in rejected)"
else bad "plan not empty after create (exit $rc)"; fi

say "== 4. secrets are not shown"
terraform show -no-color > show.txt 2>&1
grep -q -F "$PW" show.txt plan1.txt && bad "default_password visible in show/plan output" || ok "default_password hidden in show and plan"
grep -q "(sensitive value)" show.txt && ok "sensitive values masked in show" || bad "no sensitive masking found in show"

say "== 5. in-place update (alerts off, then on)"
tf apply -auto-approve -var alerts=false; result $? "update alerts_enabled=false"
terraform plan -detailed-exitcode -no-color -var alerts=false > plan2.txt 2>&1; rc=$?; cat plan2.txt >> "$LOG"
if [ $rc -eq 0 ]; then ok "no drift after update"
elif grep -q "failed to sign in" plan2.txt; then maybe "post-update drift check (sign-in rejected)"
else bad "drift after update"; fi
tf apply -auto-approve -var alerts=true; result $? "update alerts_enabled=true"

if [ "${TEST_FIREWALL:-0}" = 1 ]; then
say "== 5f. firewall"
fw_ports() { terraform state show -no-color elestio_ubuntu.t 2>/dev/null | awk -F'"' '/firewall_ports/{print $2}'; }
fw_state() { terraform state show -no-color elestio_ubuntu.t 2>/dev/null | awk '/firewall_enabled/{print $3}'; }
fw_nodrift() { # fw_nodrift <label> [terraform -var args...]
  local label="$1"; shift
  terraform plan -detailed-exitcode -no-color "$@" > planfw.txt 2>&1; local rc=$?; cat planfw.txt >> "$LOG"
  if [ $rc -eq 0 ]; then ok "$label: no drift"
  elif grep -q "failed to sign in" planfw.txt; then maybe "$label: drift check (sign-in rejected)"
  else bad "$label: drift (exit $rc)"; fi
}

say "-- default rules match what the API created (22/tcp, 4242/udp, 80/tcp, 443/tcp)"
fw_nodrift "baseline firewall"
echo "   ports now: $(fw_ports)"

say "-- add 8080/tcp"
tf apply -auto-approve -var 'fw_extra=[{port="8080",protocol="tcp"}]'; result $? "apply with extra rule 8080/tcp"
case ",$(fw_ports)," in *,8080,*) ok "API reports port 8080 after refresh" ;; *) bad "port 8080 missing from API ports: $(fw_ports)" ;; esac
fw_nodrift "after adding 8080" -var 'fw_extra=[{port="8080",protocol="tcp"}]'

say "-- replace 8080/tcp by a range 9000-9005/tcp and 8081/udp"
tf apply -auto-approve -var 'fw_extra=[{port="9000-9005",protocol="tcp"},{port="8081",protocol="udp"}]'; result $? "apply with replaced rules"
p="$(fw_ports)"
case ",$p," in *,8080,*) bad "old port 8080 still open: $p" ;; *) ok "old port 8080 removed" ;; esac
case ",$p," in *,8081,*) ok "udp 8081 present" ;; *) bad "8081 missing: $p" ;; esac
case ",$p," in *,9000-9005,*) ok "range 9000-9005 present" ;; *) bad "range 9000-9005 missing: $p" ;; esac
fw_nodrift "after replacing rules" -var 'fw_extra=[{port="9000-9005",protocol="tcp"},{port="8081",protocol="udp"}]'

say "-- disable the firewall"
tf apply -auto-approve -var fw_enabled=false; result $? "apply fw_enabled=false"
[ "$(fw_state)" = "false" ] && ok "firewall_enabled is false in state" || bad "firewall_enabled not false (got '$(fw_state)')"
fw_nodrift "firewall disabled" -var fw_enabled=false

say "-- re-enable the firewall with the default rules"
tf apply -auto-approve; result $? "apply fw_enabled=true"
[ "$(fw_state)" = "true" ] && ok "firewall_enabled is true again" || bad "firewall_enabled not true (got '$(fw_state)')"
fw_nodrift "firewall re-enabled"

say "-- invalid configuration is rejected at plan time (nothing is changed)"
terraform plan -no-color -var 'fw_extra=[{port="99999",protocol="tcp"}]' > planfw.txt 2>&1; rc=$?; cat planfw.txt >> "$LOG"
if [ $rc -ne 0 ] && ! grep -q "failed to sign in" planfw.txt; then ok "port 99999 rejected"; elif grep -q "failed to sign in" planfw.txt; then maybe "port 99999 check (sign-in rejected)"; else bad "port 99999 was accepted"; fi
terraform plan -no-color -var 'fw_extra=[{port="8080",protocol="icmp"}]' > planfw.txt 2>&1; rc=$?; cat planfw.txt >> "$LOG"
if [ $rc -ne 0 ] && ! grep -q "failed to sign in" planfw.txt; then ok "protocol icmp rejected"; elif grep -q "failed to sign in" planfw.txt; then maybe "icmp check (sign-in rejected)"; else bad "protocol icmp was accepted"; fi
terraform plan -no-color -var fw_drop_ssh=true > planfw.txt 2>&1; rc=$?; cat planfw.txt >> "$LOG"
if [ $rc -ne 0 ] && grep -q "22" planfw.txt && ! grep -q "failed to sign in" planfw.txt; then ok "dropping required port 22/tcp (SSH) rejected"; elif grep -q "failed to sign in" planfw.txt; then maybe "ssh-port check (sign-in rejected)"; else bad "removing 22/tcp was accepted"; fi
fi

say "== 6. refresh keeps a healthy service"
tf apply -refresh-only -auto-approve; result $? "refresh-only apply"
terraform state list 2>/dev/null | grep -q elestio_ubuntu.t && ok "resource still in state" || bad "healthy resource dropped from state"

say "== 7. delete OUTSIDE Terraform (waits up to 10 min for the backend)"
out=$(liveapi delete_and_wait "$PROJECT" "$SVC_ID"); rc=$?
echo "$out" | sed 's/^/   /' | tee -a "$REPORT"
if [ $rc -eq 2 ]; then maybe "out-of-band delete (sign-in rejected)"
elif echo "$out" | grep -q "GONE yes"; then ok "backend answers InvalidServer/service_deleted for the gone service"
else bad "out-of-band delete or wait failed"; fi

say "== 8. refresh must drop the missing service from state"
terraform plan -no-color > plan3.txt 2>&1; cat plan3.txt >> "$LOG"
if grep -q "failed to sign in" plan3.txt; then maybe "missing-service detection (sign-in rejected)"
elif grep -q "1 to add" plan3.txt; then ok "refresh removed it from state; plan wants to re-create"
else bad "plan did not detect the missing service"; fi

say "== 9. destroy of an already-deleted service must succeed"
tf destroy -auto-approve; rc=$?; result $rc "destroy succeeded"
if terraform state list 2>/dev/null | grep -q elestio_ubuntu; then
  [ $rc -ne 0 ] && auth_failed && maybe "state still has the resource (destroy was rejected at sign-in)" || bad "resource still in state"
else ok "state is empty"; fi

say "== 10. ordinary create + destroy (delete wait must recognize a real deletion)"
# The backend keeps the name of a service in 'deleting' state reserved, so use a new name.
NAME="tf-audit-$(date +%H%M%S)b"
sed -i "s/server_name      = \".*\"/server_name      = \"$NAME\"/" main.tf
tf apply -auto-approve; result $? "second create succeeded"
tf destroy -auto-approve; result $? "normal destroy succeeded"

say "== 11. final check: nothing left in project $PROJECT"
out=$(liveapi count "$PROJECT" "tf-audit-"); rc=$?
say "   $out"
if [ $rc -eq 2 ]; then maybe "leftover check (sign-in rejected) - CHECK THE DASHBOARD FOR tf-audit-*"
elif [ "$out" = "COUNT 0" ]; then ok "no active tf-audit-* service remains in project $PROJECT"
else bad "an active tf-audit-* service STILL EXISTS - delete it in the dashboard"; fi

say ""
finish
say ""
say "RESULT: $pass passed, $fail failed, $inconc inconclusive"
say "Send: $REPORT  and  $OUT/terraform$SUFFIX.redacted.log"
[ "$fail" = 0 ] && [ "$inconc" = 0 ]
