# miniasm

`miniasm` is a simple cross-platform Assembly, which lets you write programs in Assembly fast and easy:
```miniasm
main { 
  // Initialize variables
  mov x, 10; 
  mov y, 20; 
  mov result, 0;

  // Move x+y into result
  add result, x, y;  

  // Print result into console with formatting
  print "{}", result; 
  ret 0; 
}
```
