START_PROGRAM
decr={
  _2=-1;
  _r=add();
}

loop={
  _1=stop;
  stop = decr();

  payload();

  _1=stop;
  _2=loop;
  if();
}

stop=10;
payload={
  _1=stop;
  print();
}

loop();

END_PROGRAM
START_OUTPUT
9
8
7
6
5
4
3
2
1
0
END_OUTPUT
