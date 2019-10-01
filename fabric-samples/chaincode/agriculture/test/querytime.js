const exec = require('child_process').exec
const spawn = require('child_process').spawn

console.time('test1')

var child = exec('sh querytime.sh', function(error, stdout, stderr){
  //console.info('stdout: ')
  //console.log(stdout)
})

console.timeEnd('test1')

//peer chaincode query -C $CHANNEL_NAME -n mycc -c '{"Args":["queryProduct", "No998"]}'

//let ps = spawn('ant.bat', ['-propertyfile ', prop1, '-propertyfile ', prop2, '-Dgit.branch='+branch, '-Dinstallfile.suffix="'+prop4+'"'], { shell:true });