LD r5 data1
LD r6 data3
LD r7 data5
ADD R0 R5 r6
ADD R0 R0 R7
HALT
data1: .FILL 3
data2: .FILL 4
data3: .FILL 5
data4: .FILL 6
data5: .FILL 7