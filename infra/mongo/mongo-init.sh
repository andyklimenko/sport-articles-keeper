set -e

mongo <<EOF

use ${MONGO_INITDB_DATABASE}
db.createCollection('newsletterNews')
db.newsletterNews.createIndex( { newsArticleId: 1} )
db.newsletterNews.createIndex( { id: 1} )

EOF