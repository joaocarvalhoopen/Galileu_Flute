Game title: Galileu's Flute

Description: This is a game to be played with a real Recorder (Flauta de Bisel
             or Flauta Doce).
             The game analyses the audio of the flute and determines
             what musical note was played, it shows a musical score
             that is scrolling from the left to the right and the
             objective is to hit each note at the rigth time.
             This game is in text mode and is tested on Windows 10.

Author:  Joao Nuno Carvalho
email:   joaonunocarv@gmail.com
Date:    2017-12-08
License: MIT OpenSource license


To execute do on a command line:
      galileu_flute.exe
   or reading from a specific json file
      galileu_flute.exe ./music_02.json
   or reading from a specific simplified ABC file format (extension .ABC or .abc)
      galileu_flute.exe ./music_01.ABC


Example of output:

         Galileu's Flute

             Score: 156
  ---
 | = |
 |   |
 | # | #  |.............................
 | O |    |.........S_D.................
 | O |    |.............................
 | O |    @_......S_....................
 | O |    |.S_..........................
 | O |    |...S_........................
 |O  |    |.....S_......................
  | |
 |   |
  ---

 At the Flute:
   'O' - an open hole, no finger.
   '#' - an closed hole, put your finger.

 In the Sheet Music:
   'S' - a normal note, from DO to SI.
   '_' - the continuation of the same note.
   'D' - the DO_HIGH.
   '.' - an indication of no note.

 At the Score Line:
   'X' - You hit the correct note, 10 point's.
   '@' - You hit the wrong note, -1 point.
   '|' - Just the indication of the line.

  To Exit the program:
    Wait for 15 minuts or hit Ctrl + c keys.

 Note: The note detection algorithm is based on the frequency detection
 algoritm of YIN'ss algorithm more specifically the implementationm on
 https://github.com/ashokfernandez/Yin-Pitch-Tracking/blob/master/Yin.c

Codes to create or modify the music_01.json file:

Code for each note:
	EMPTY    = 0
	DO       = 1
	RE       = 2
	MI       = 3
	FA       = 4
	SOL      = 5
	LA       = 6
	SI       = 7
	DO_HIGH  = 8


Simplified ABC file format ( *.ABC or *.abc ):

  It starts with a
  T: <Name of the music>
  Followed by the note
  C, D, E, F, G, A, B, c
  This are the only notes that are recognizes by the game.
  After the note it can have a number that is optional that is
  1, 2, 3, 4 and marks the duration of the note.
  Example C1,C2, C3, C4, D1, E4, c4
  Silences can be made with the note S, and can also be followed
  by a number for the duration.
  The symbol space, |, or ] are ignored.
  See the example file for an example.
  With this file format is easy to transform the song written in
  the ABC format into this simplified ABC format.

Example of the simplified ABC file format:

T: Easy Music
C D E F| G A B c|]
CDEF|GABc|S2C2]
C1D1E1F1|G1A1B1c1|]
C2D2E2F2|G2A2B2c2|]
C3D3E3F3|G3A3B3c3|]
C4D4E4F4|G4A4B4c4|]

