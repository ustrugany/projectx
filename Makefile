arg = $(filter-out $@,$(MAKECMDGOALS))

.PHONY: help
list:
	@$(MAKE) -qp | awk -F':' '/^[a-zA-Z0-9][^$$#\/\t=]*:([^=]|$$)/ {split($$1,A,/ /);for(i in A)print A[i]}' | grep -v Makefile | sort -u

help::
	@echo "Makefile"

include Makefile.*
-include pf-makefiles/Makefile.*

%:
	@: