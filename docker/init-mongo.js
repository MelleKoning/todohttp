db.createUser(
    {
        user : "tododbuser",
        pwd : "tododbpass",
        roles : [
            {
                role: "readWrite",
                db: "mongotododb"
            }
        ]
    }
)