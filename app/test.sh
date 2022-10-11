go build

echo "Test1:"
./app tests/file1.txt

echo "Test2:"
./app tests/file2.txt

echo "Test3:"
./app tests/file1.txt tests/file2.txt

echo "Test4:"
cat tests/file1.txt |./app
