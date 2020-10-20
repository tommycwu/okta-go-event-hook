// server.js
// where your node app starts


// Current project dependencies, feel free to add additional libraries or frameworks in `package.json`.
const express = require("express");
const bodyParser = require('body-parser');
const app_request = require('request');
const fs = require("fs");
const app = express();

// parse application/x-www-form-urlencoded
app.use(bodyParser.urlencoded({ extended: false }))

// parse application/json
app.use(bodyParser.json())


// init sqlite db
const dbFile = "./.data/sqlite.db";
const exists = fs.existsSync(dbFile);
const sqlite3 = require("sqlite3").verbose();
const db = new sqlite3.Database(dbFile);


// if ./.data/sqlite.db does not exist, create it, otherwise print records to console
db.serialize(() => {
  if (!exists) {
    db.run(
      "CREATE TABLE Users (id INTEGER PRIMARY KEY AUTOINCREMENT, firstName TEXT, lastName TEXT, email TEXT, mobilePhone Int, oktaId Int)"
    );
    console.log("New table Users created!");
    // insert first row of User Data
    db.serialize(() => {
      db.run(
        'INSERT INTO Users (firstName, lastName, email, mobilePhone, oktaId) VALUES ("Test", "User", "testuser@doesnotexit.com", "111-222-3333", 1234)'
      );
    });
  }
});


// Event Hook Verification
// Extract header 'x-okta-verification-challenge' from Okta request
// Return value as JSON object verification
app.get("/userTransfer", (request, response) => {
  var returnValue = {
    "verification": request.headers['x-okta-verification-challenge'],
  };  
  response.json(returnValue);
});


// transferContact request, POST from Okta
app.post("/userTransfer", (request, response) => {
  var newUser = request.body.data.events[0]['target'][0];
  response.sendStatus(200);
  dbInsert(newUser);
  console.log('New User Inserted');
  }
)
  
// Request additional data from Okta using getUserData function
// Push new user to database
async function dbInsert(newUser) {
  var profile = await getUserData(newUser['alternateId']);
  var query =  "INSERT INTO Users (firstName, lastName, email, mobilePhone, oktaId) VALUES ("
  query += '"' + profile['firstName'] + '", "' + profile['lastName'] + '", "' + profile['email'] + '", "' + profile['mobilePhone'] + '", "' + newUser['id'] + '")'
  db.run(
    query
  ); 
}

// Request additional data from Okta
// Return data to dbInsert function
function getUserData(newUser) {
  return new Promise(resolve => {
    var options = {
      'method': 'GET',
      'url': process.env.okta_url + '/api/v1/users?q=' + newUser,
      'headers': {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': 'SSWS ' + process.env.okta_key
      }
    };

    app_request(options, function (error, response) { 
      if (error) throw new Error(error);
      var result = JSON.parse(response.body)[0]['profile']; 
      resolve(result);
    });
  })
}


// Upon request, return all users in database
app.get("/getUsers", async (request, response) => {
  let promise = new Promise((res, rej) => {
      db.all("SELECT * from Users", (err, rows) => {
        res(rows);
  })});
  
  var users = await promise;
  response.send(JSON.stringify(users));

});


// make all the files in 'public' available
// https://expressjs.com/en/starter/static-files.html
app.use(express.static("public"));

// https://expressjs.com/en/starter/basic-routing.html
app.get("/", (request, response) => {
  response.sendFile(__dirname + "/views/index.html");
});



// listen for requests :)
const listener = app.listen(process.env.PORT, () => {
  console.log("Your app is listening on port " + listener.address().port);
});
