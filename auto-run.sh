pushd python
virtualenv env --python=python3
source env/bin/activate
pip install -r requirements.txt
python server.py &
popd

pushd go
go build -o run .
./run
popd