.ORIG x3000
ADD R7 R7 #0 ; Must do at least one operation to set the COND register to allow BR to work
BR foo
ADD R0 R0 #1
HALT
foo: ADD R0 R0 #2
HALT
bar: ADD R0 R0 #3
HALT