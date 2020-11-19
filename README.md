# ascii
##Description:
Transforms source image (png/jpeg) to ascii representation

##Build:
go install -v ...

##Usage:
 Example: 
 
 `ascii -in "/tmp/fish.png" -out "/tmp/ascii.txt" -p -w 135 -h 52 -c`   
 
 Flags:
 ```
  -c	compute concurrently or not
  -h int
    	output pic height in chars (default 50)
  -in string
    	path to image file (jpeg/png) (default "example.png")
  -out string
    	output txt file (default "ascii.txt")
  -p	print the result to console or not
  -w int
    	output pic width in chars (default 80)
 ```