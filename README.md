# miniasm

`miniasm` is a simple cross-platform Assembly, which lets you write programs in Assembly fast and easy:
```miniasm
main (args []) { 
  mov x, 10; 
  mov y, 20; 
  mov result, 0;

  add result, x, y;  

  print "{}", result; 
  ret 0; 
}
```
