#!/bin/bash

# Check parameters
case $# in
   2)
      env=$1
      port=$2
   ;;
    *)
      echo "Usage: sh $0 test 5601"
      exit 1;
    ;;
esac

# Check env
if [ "${env}" != "test" ] && [ "${env}" != "uat" ] && [ "${env}" != "dev" ];then
   echo "env is error.not support ${env}"
   exit 1;
fi

# get base dir
baseDir=$(cd `dirname $0`;cd ..;pwd)
cd ${baseDir}

projectName="`basename ${baseDir}`_${env}"
serverName="${projectName}_server"
buildDir="${baseDir}/src/server"
binDir="${baseDir}/src/bin"
configDir="${baseDir}/src/config"
releaseDir="${baseDir}/release"

pkgServer(){
   [ -d "${releaseDir}/${projectName}" ] && rm -rf ${releaseDir}/${projectName}
   mkdir -p ${releaseDir}/${projectName}
   [ -d "${releaseDir}/${projectName}/config" ] || mkdir -p ${releaseDir}/${projectName}/config
   [ -d "${releaseDir}/${projectName}/log" ] || mkdir -p ${releaseDir}/${projectName}/log
   [ -d "${releaseDir}/${projectName}/run" ] || mkdir -p ${releaseDir}/${projectName}/run
   [ -d "${releaseDir}/${projectName}/bin" ] || mkdir -p ${releaseDir}/${projectName}/bin

   cp -rf ${configDir}/* ${releaseDir}/${projectName}/config

   awk '{
      if($1~"port=") {
         $1="port='${port}'";
      } else if($1~"serverName=\"") {
         $1="serverName=\"'${serverName}'\""
      } else if($1~"env=\"") {
         $1="env=\"'${env}'\""
      }
      print
   }' ${binDir}/start.sh > ${releaseDir}/${projectName}/bin/start.sh

   cd ${buildDir} && go build -o ${releaseDir}/${projectName}/bin/${serverName}

   cd ${releaseDir}
   [ -f "${projectName}.latest.tar.gz" ] && mv ${projectName}.latest.tar.gz ${projectName}.`date +"%Y-%m-%d_%H-%M-%S"`.tar.gz
   tar czvf ${projectName}.latest.tar.gz ${projectName} >/dev/null 2>&1
   rm -rf ${projectName}
   cd ${baseDir}
   echo "Packaging success.See file ${releaseDir}/${projectName}.latest.tar.gz."
}

pkgServer