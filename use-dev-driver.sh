# This script is is not designed to be invoked directly; it is designed
# to used via "source ./use-dev-driver.sh" or ". ./use-dev-driver.sh"

make dev
export PATH=$PWD/_bin:$PATH
