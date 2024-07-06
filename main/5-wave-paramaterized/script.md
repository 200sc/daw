# script 

There are a couple of assumptions that we've been accidentally running with that some more experimentation
will show us we need to adjust.

One assumption; our speaker is only one speaker. If I'm wearing headphones, all sorts of audio can positionally or otherwise
decide that some sounds will come out of one of my headphone ears, but not the other, or be quieter from one or the other.

To demonstrate this requires some brute force, but eliding some experiments we will find that the first four bytes of our
input go to one of our headphones and the second four go to the other.

Which reveals our second problem: we've been treating this speaker as if it cares about individual bytes. Apparently, it cares about
sets of four.

And that brings us to a bigger example. Lets encode our triangle wave into int32 values instead of bytes. And let's encode that same value for
as many channels as we have. 

- run 

Note the 128 value here, this reflects our period or frequency. 

- run with 64

and note with half of that 128 value, the pitch sounds similar if lower than 128

- run with 32

and again with 32, it sounds similar

- run with 39 

but with 39 it does not. What we're discovering has a few names, resonance or more commonly in western music, octaves. 
In western music a given pitch sounds the same as many other pitches, i.e C4 vs C5. These pitches have equivalent freqencies,
but multiplied or divided by two. 

t: 3m