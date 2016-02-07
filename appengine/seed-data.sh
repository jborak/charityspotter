

BASE_URL=https://charitysandbox.firebaseio.com

curl -X POST -d '{ "created": 1446998421, "tags": [ "boots", "shoes", "brown" ], "uid": "1", "url": "/demo/brown1.jpg" }' $BASE_URL/items.json
curl -X POST -d '{ "created": 1446998421, "tags": [ "boots", "shoes", "brown" ], "uid": "2", "url": "/demo/brown2.jpg" }' $BASE_URL/items.json
curl -X POST -d '{ "created": 1446998421, "tags": [ "boots", "shoes", "black" ], "uid": "3", "url": "/demo/black1.jpg" }' $BASE_URL/items.json
curl -X POST -d '{ "created": 1446998421, "tags": [ "boots", "shoes", "black" ], "uid": "4", "url": "/demo/black2.jpg" }' $BASE_URL/items.json
curl -X POST -d '{ "created": 1446998421, "tags": [ "pants", "yellow" ], "uid": "5", "url": "/demo/yellow1.jpg" }' $BASE_URL/items.json
curl -X POST -d '{ "created": 1446998421, "tags": [ "pants", "gray", "grey" ], "uid": "6", "url": "/demo/gray1.jpg" }' $BASE_URL/items.json
curl -X POST -d '{ "created": 1446998421, "tags": [ "pants", "tan" ], "uid": "7", "url": "/demo/tan1.jpg"} ' $BASE_URL/items.json
