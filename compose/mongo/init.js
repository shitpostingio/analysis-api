db = db.getSiblingDB('analysis')

db.createUser(
    {
        user: "analysis",
        pwd: "analysis",
        roles: [
            {
                role: "readWrite",
                db: "analysis"
            }
        ]

    }
);

db.tokens.insertOne({"token":<TOKEN-HERE>});