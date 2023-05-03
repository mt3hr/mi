$wd = pwd
rm -Recurse -Force .\src\app\mi\mi\embed\html\*
npm run build
cp -Recurse .\dist\* .\src\app\mi\mi\embed\html\

cd .\src\app\mi
go install .
cd $wd
echo "mi install finished!"