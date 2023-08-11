set -e

mongo <<EOF

use ${MONGO_INITDB_DATABASE}
db.createCollection('newsletterNews')

EOF