export FDK_AAC_VERSION="2.0.0"
export CGO_CPPFLAGS="-I/usr/local/fdk-aac-${FDK_AAC_VERSION}/include/fdk-aac"
export CGO_LDFLAGS="-L/usr/local/fdk-aac-${FDK_AAC_VERSION}/lib"
export LD_LIBRARY_PATH="/usr/local/fdk-aac-${FDK_AAC_VERSION}/lib"
export LOG_LEVEL=debug