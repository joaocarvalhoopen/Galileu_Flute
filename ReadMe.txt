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
    or
   galileu_flute.exe ./music_02.json


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
