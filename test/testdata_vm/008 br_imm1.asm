.ORIG x3000
ADD R7 R7 #1 ; Must do at least one operation to set the COND register to allow BR to work
BR #1
ADD R0 R0 #1 ; this line will be skipped
ADD R1 R1 #1
BR #0
ADD R2 R2 #1
HALT
ADD R3 R3 #1