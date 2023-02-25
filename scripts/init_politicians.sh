curl --location 'localhost:8080/politician' \
--header 'Content-Type: application/json' \
--data '{
    "name": "王婉諭",
    "birthdate": "1979-04-26",
    "avatarUrl": "https://upload.wikimedia.org/wikipedia/commons/4/48/%E7%AB%8B%E6%B3%95%E5%A7%94%E5%93%A1%E7%8E%8B%E5%A9%89%E8%AB%AD.jpg"
}'

curl --location 'localhost:8080/politician' \
--header 'Content-Type: application/json' \
--data '{
    "name": "王世堅",
    "birthdate": "1960-01-01",
    "avatarUrl": "https://www.ly.gov.tw/Images/Legislators/ly1000_6_00003_23f.jpg"
}'

curl --location 'localhost:8080/politician' \
--header 'Content-Type: application/json' \
--data '{
    "name": "許淑華",
    "birthdate": "1975-10-15",
    "avatarUrl": "https://upload.wikimedia.org/wikipedia/commons/e/e7/%E8%A8%B1.JPG"
}'

curl --location 'localhost:8080/politician' \
--header 'Content-Type: application/json' \
--data '{
    "name": "許淑華",
    "birthdate": "1973-05-22",
    "avatarUrl": "https://upload.wikimedia.org/wikipedia/commons/c/c4/Hsu_Shu-Hua_at_World_Design_Capital_Taipei_press_conference_20120629.jpg"
}'
