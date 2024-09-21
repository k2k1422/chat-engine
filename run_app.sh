set -a
DEFAULT_FILENAME=".env"
source $DEFAULT_FILENAME
FILENAME="${1:-$DEFAULT_FILENAME}"
source $FILENAME
set +a
go run main.go