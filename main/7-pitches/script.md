# script 

enter, math. We now know that frequencies correspond to some pitch, but up until now have been approximating exactly how
we represent these frequencies. Lets use some known values for frequencies.

- go to daw.A7

reflected here is some version of western music's understanding of canonical pitches. A through G, flats and sharps.
And we can see here as we found in our experiments, B7 is half of B8, and etc. These values represnet how often a pattern
should repeat in order to output a specific pitch. How do we know these values? We loop em up.

How can we take a frequency and make our waveforms reflect it? Well, if we had a function that told us how far we were
into a pattern, and our waveforms could output parts of the pattern that reflect that percentage, we'd be good, right?

So we take our current sample, divide by how many samples we should process per second, and multiply by the hertz of our
pitch. We also put this on a scale of 2 pi so it can be consumed by the sin function here in a clean loop.

- run

And we get a clean sin wave at our specified pitch.

We've mulled over a few things; sample rate, format, and here in our goroutine we're calling WritePCM. Well, we've
reached the point where we know enough fundamentals that we can reveal the contract our speaker expects us to use.

Our writer expects us to send it 44.1k samples every second, on two channels, and it expects each sample to have 32 bits
of data. We can actually change any of these values, but these are reasonable defaults. PCM or Pulse Code Modulation is
the name for this format, defined in terms of sample rate and bit rate. PCM is the standard for digital audio.

t: 3m