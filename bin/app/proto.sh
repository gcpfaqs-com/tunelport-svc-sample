echo 
echo "# ============================================="
echo "#  __  __  ____  _  __"
echo "# |  \\/  |  _ \\| |/ /"
echo "# | \\  / | | | | ' / "
echo "# | |\\/| | | | |  <  "
echo "# | |  | | |_| | . \\ "
echo "# |_|  |_|____/|_|\\_\\"
echo "#                                   "
echo "#  Creador: Marco Antonio - markitos      "
echo "# ============================================="
echo "#  🥷 (mArKit0sDevSecOpsKit) 🗡️"
echo "#  Markitos DevSecOps Kulture"
echo "# ============================================="
echo 

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
cd "$REPO_ROOT"

set -euo pipefail
IFS=$'\n\t'

function log_info() {
    echo "[INFO] $*"
}

function log_error() {
    echo "[ERROR] $*" >&2
}

function require_cmd() {
    if ! command -v "$1" >/dev/null 2>&1; then
        log_error "No se encontro el comando requerido en PATH: $1"
        exit 1
    fi
}

PROTO_DIR="internal/infrastructure/proto"
OUT_DIR="internal/infrastructure/grpc"

log_info "Directorio de trabajo: $REPO_ROOT"
log_info "Argumentos proporcionados: ${*:-ninguno}"

require_cmd protoc
require_cmd protoc-gen-go
require_cmd protoc-gen-go-grpc

if [ ! -d "$PROTO_DIR" ]; then
    log_error "No existe el directorio de protos: $PROTO_DIR"
    exit 1
fi

mapfile -t proto_files < <(find "$PROTO_DIR" -type f -name "*.proto" | sort)

if [ "${#proto_files[@]}" -eq 0 ]; then
    log_error "No se encontraron archivos .proto en $PROTO_DIR"
    exit 1
fi

mkdir -p "$OUT_DIR"
rm -f "$OUT_DIR"/*.pb.go

log_info "Compilando ${#proto_files[@]} archivo(s) .proto..."

protoc \
    --proto_path="$PROTO_DIR" \
    --go_out="$OUT_DIR" \
    --go-grpc_out="$OUT_DIR" \
    --go_opt=paths=source_relative \
    --go-grpc_opt=paths=source_relative \
    "${proto_files[@]}"

log_info "Generacion completada en $OUT_DIR:"
find "$OUT_DIR" -type f -name "*.pb.go" | sort