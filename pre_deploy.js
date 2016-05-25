/// CMD PROD

var fs=require("fs");

const DEV_MAIN="main.go";
const PROD_MAIN="cmd/smartdb/main.go";

//GO Deps
const GODEP="Godeps/Godeps.json";
const INTERNAL="github.com/Zombispormedio/smartdb/";
const PACKAGES=[
  "routes", "config", "models" , "controllers", "middleware", "lib/response", "lib/utils", "lib/struts", "lib/store", "lib/rabbit",
   "consumer"
];

fs.writeFileSync(PROD_MAIN, fs.readFileSync(DEV_MAIN));



var pre_godep=JSON.parse(fs.readFileSync(GODEP));

pre_godep.Deps=pre_godep.Deps.concat(
    PACKAGES.map(function(pkg){return {"ImportPath":INTERNAL+pkg};})
);

fs.writeFileSync(GODEP, JSON.stringify(pre_godep));