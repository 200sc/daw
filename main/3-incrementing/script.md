# script

To confirm our theory, lets write incrementing values and see what happens. 

And, lets start displaying exactly what we're writing to this writer.

- run

That's a little interesting; it'd be nice if this kept going for longer though 
-- loop over 100000

- run

now we've got continuous sound, and we've got a variable we can manipulate, so lets manipulate it

-- writeIncrementing 10000 size

- run

-- writeIncrementing 5000 size

- run

-- writeIncrementing 2500 size

- run

-- writeIncrementing 1250 size

- run

We've discovered something new-- pitch. This size value we've written here is alternately, our period-- how often a pattern
we've chosen is being written to the audio interface. The faster that pattern is written, the higher pitch it outputs.  


1m30