#! /bin/bash -

if [[ $1 == "all" ]] ; then
	files_path=`ls *.go`
else 
	files_path=`ls -t *.go | sed 1q`
fi

echo "*************************************************"
echo "---------START-------------"
echo "ALL FILES: "$files_path

#ls *.go | while read line
while read line 
do
	echo "----------------------"
	echo "go build $line waiting."
	gofmt -w $line
	go build $line
	[[ $? -ne 0 ]] && echo "go build $line fait." && exit 1
	echo "go build $line success."
#	echo "----------------------"
done <<< "$files_path"

echo "----------------------"
run_file=`echo $files_path | awk -F "." '{print $1}'`
echo "*************************************************"
echo "-------RUN-- $run_file ---START----------"
time ./$run_file
echo "----------------------"
echo "-------RUN-- $run_file ----END-----------"
echo "----------END--------------"
echo "*************************************************"

./test.sh
