.PHONY: all build maketestdata
.DEFAULT_GOAL := build
DEST := 'deepfindexe'

build:
	@-GOOS=darwin GOARCH=amd64 go build -o bin/Darwin/${DEST} ./cmd/main.go
	@-GOOS=linux GOARCH=amd64 go build -o bin/Linux/${DEST} ./cmd/main.go
	@-GOOS=windows GOARCH=amd64 go build -o bin/Windows/${DEST}.exe ./cmd/main.go

# for regenerate test data
# brew install brotli lz4 xz snzip zstd p7zip
# brew install --cask rar
maketestdata: clean
	@-cd testdata && \
	echo "test bat" > test.bat && \
	echo "\x4D\x5A_test_mime_exec" > test.exe
	@-echo prepare stub bat and exe

	@-cd testdata && \
	tar -cf test.bat.tar test.bat && \
	zip -q test.bat.zip test.bat  && \
	gzip --keep test.bat && \
	bzip2 --keep test.bat && \
	lz4 -q test.bat && \
	xz --keep test.bat && \
	snzip -k test.bat && \
	zstd -q test.bat && \
	brotli test.bat && \
	rar a -inul test.bat.rar test.bat && \
	7z a test.bat.7z test.bat > /dev/null
	@-echo bat single level compression

	@-cd testdata && \
	tar -cf test.exe.tar test.exe && \
	zip -q test.exe.zip test.exe  && \
	gzip --keep test.exe && \
	bzip2 --keep test.exe && \
	lz4 -q test.exe && \
	xz --keep test.exe && \
	snzip -k test.exe && \
	zstd -q test.exe && \
	brotli test.exe && \
    rar a -inul test.exe.rar test.exe && \
    7z a test.exe.7z test.exe > /dev/null
	@-echo exe single level compression

	@-cd testdata && \
	zip -q test.bat.tar.zip test.bat.tar && \
	gzip --keep test.bat.tar && \
	bzip2 --keep test.bat.tar && \
	lz4 -q test.bat.tar && \
	xz --keep test.bat.tar && \
	snzip -k test.bat.tar && \
	zstd -q test.bat.tar && \
	brotli test.bat.tar && \
    rar a -inul test.bat.tar.rar test.bat.tar && \
    7z a test.bat.tar.7z test.bat.tar > /dev/null
	@-echo bat with tar

	@-cd testdata && \
	zip -q test.exe.tar.zip test.exe.tar && \
	gzip --keep test.exe.tar && \
	bzip2 --keep test.exe.tar && \
	lz4 -q test.exe.tar && \
	xz --keep test.exe.tar && \
	snzip -k test.exe.tar && \
	zstd -q test.exe.tar && \
	brotli test.exe.tar && \
    rar a -inul test.exe.tar.rar test.exe.tar && \
    7z a test.exe.tar.rar.7z test.exe.tar > /dev/null
	@-echo exe with tar

	@-cd testdata && \
	zip -q test.exe.recursive.zip test.exe.rar && \
	rar a -inul test.exe.recursive.rar test.exe.zip && \
    7z a test.exe.recursive.7z test.exe.zip > /dev/null
	@-echo exe second level

clean:
	rm -f testdata/test*

test:
	go test
