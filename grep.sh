grep --line-number --exclude-dir={vendor,.git,.gitsecret,fontawesome,img} --exclude={Gopkg.*,\*.log} -r -e $1
