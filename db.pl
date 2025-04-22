switch(x).
switch(w).
switch(z).
seen(x,d,1).
seen(x,w,8).
seen(x,c,8).
seen(x,a,8).
seen(x,b,8).
seen(x,z,8).
seen(w,c,6).
seen(w,b,7).
seen(w,x,10).
seen(w,d,10).
seen(w,a,4).
seen(w,z,4).
seen(w,b,4).
seen(z,a,1).
seen(z,b,2).
seen(z,c,5).
seen(z,d,5).
seen(z,w,5).
seen(z,x,5).

far(X,Y) :-	seen(X,K,MIDPORTX), 
		seen(Y,K,MIDPORTY),
		seen(X,Y,PORTX),
		seen(Y,X,PORTY),
		switch(K), switch(X), switch(Y), X \= Y, Y \= K, X \= K, 
		PORTX == MIDPORTX,
		PORTY == MIDPORTY.

direct(X,Y) :-	seen(X,K,MIDPORTX),
		seen(Y,K,MIDPORTY),
		seen(X,Y,PORTX),
		seen(Y,X,PORTY),
		switch(K), switch(X), switch(Y), X \= Y, Y \= K, X \= K,
		(PORTX \= MIDPORTX ; PORTY \= MIDPORTY).