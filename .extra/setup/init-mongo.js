function seed(dbName, user, password) {
  db = db.getSiblingDB(dbName);
  db.createUser({
    user: user,
    pwd: password,
    roles: [{ role: "readWrite", db: dbName }],
  });

  db.createCollection("api_keys");
  db.createCollection("roles");

  db.api_keys.insert({
    key: "1D3F2DD1A5DE725DD4DF1D82BBB37",
    permissions: ["GENERAL"],
    comments: ["To be used by the xyz vendor"],
    version: 1,
    status: true,
    createdAt: new Date(),
    updatedAt: new Date(),
  });

  db.roles.insertMany([
    {
      code: "LEARNER",
      status: true,
      createdAt: new Date(),
      updatedAt: new Date(),
    },
    {
      code: "WRITER",
      status: true,
      createdAt: new Date(),
      updatedAt: new Date(),
    },
    {
      code: "EDITOR",
      status: true,
      createdAt: new Date(),
      updatedAt: new Date(),
    },
    {
      code: "ADMIN",
      status: true,
      createdAt: new Date(),
      updatedAt: new Date(),
    },
  ]);

  db.users.insert({
    name: "Admin",
    email: "admin@unusualcode.org",
    password: "$2a$10$psWmSrmtyZYvtIt/FuJL1OLqsK3iR1fZz5.wUYFuSNkkt.EOX9mLa", // hash of password: changeit
    roles: db.roles
      .find({})
      .toArray()
      .map((role) => role._id),
    status: true,
    createdAt: new Date(),
    updatedAt: new Date(),
  });
}

seed("unusualcode-org-prod-db", "unusualcode-org-prod-db-user", "changeit");
seed("unusualcode-org-dev-db", "unusualcode-org-dev-db-user", "changeit");
seed("unusualcode-org-test-db", "unusualcode-org-test-db-user", "changeit");
