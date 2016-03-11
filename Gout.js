var child_process=require('child_process');
var exec = child_process.exec;
var spawn = child_process.spawn;

const GO_INSTALL = 'go install';
const SERVICE='smartdb.exe';


function run_cmd(cmd, args, cb, error, end) {
    child = spawn(cmd, args),
        me=this;

    child.stdout.on('data', function (buffer) { cb(me, buffer); });
    child.stderr.on('data', function(data) {
        error(data);
    });
    child.stdout.on('end', end);
    return child;
}

function execMain(){
    console.log("Starting service...");
    var proc=run_cmd(SERVICE, [], function(me, buffer){
        console.log(buffer.toString());
    },function(data){
        console.log(data.toString());

    }, function(){
        console.log("Service finished");
         install();
    });


    var stdin = process.openStdin();
    stdin.addListener("data", function(d) {
        var command=d.toString().trim();

        switch(command){
            case "rs":{
                console.log("Killing service");
                proc.kill('SIGINT');
                break;
            }
            case "stop": {
                console.log("Process exit()");
                process.exit();
                break;
            }

        }

    });

}



function install(){
    console.log("Compiling service...")
    exec(GO_INSTALL, function(error, stdout, stderr) {
        if(error){
            console.log(error);
            process.exit();
        } 
        console.log("Service compiled");
        execMain(); 
    });
}

install();
