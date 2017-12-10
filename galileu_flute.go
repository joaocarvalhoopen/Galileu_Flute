// Game title: Galileu's Flute
//
// Description: This is a game to be played with a real Recorder (Flauta de Bisel
//              or Flauta Doce).
//              The game analyses the audio of the flute and determines
//              what musical note was played, it shows a musical score
//              that is scrolling from the left to the right and the
//              objective is to hit each note at the rigth time.
//              This game is in text mode and is tested on Windows 10,
//              but it can run in Linux or Mac.
//              Implemented in Go ( GoLang ) programming Language.
//
//
// Author:  Joao Nuno Carvalho
// email:   joaonunocarv@gmail.com
// Date:    2017-12-08
// License: MIT OpenSource license
//
// To execute do on a command line:
//    galileu_flute.exe
//    or
//    galileu_flute.exe ./music_02.json
//
// Example of the output:
//
//
//         Galileu's Flute
//
//              Score: 156
//   ---
//  | = |
//  |   |
//  | # | #  |.............................
//  | O |    |.........S_D.................
//  | O |    |.............................
//  | O |    @_......S_....................
//  | O |    |.S_..........................
//  | O |    |...S_........................
//  |O  |    |.....S_......................
//   | |
//  |   |
//   ---
//
//  At the Flute:
//    'O' - an open hole, no finger.
//    '#' - an closed hole, put your finger.
//
//  In the Sheet Music:
//	  'S' - a normal note, from DO to SI.
//    '_' - the continuation of the same note.
//    'D' - the DO_HIGH.
//    '.' - an indication of no note.
//
//  At the Score Line:
//   'X' - You hit the correct note, 10 point's.
//   '@' - You hit the wrong note, -1 point.
//   '|' - Just the indication of the line.
//
//  To Exit the program:
//     Wait for 15 minuts or hit Ctrl + c keys.
//
// Note: The note detection algorithm is based on the frequency detection
// algoritm of YIN'ss algorithm more specifically the implementationm on
// https://github.com/ashokfernandez/Yin-Pitch-Tracking/blob/master/Yin.c
//

package main

import (
	"github.com/gordonklaus/portaudio"
	"time"
	"fmt"
	"math"
	"bytes"
	"io/ioutil"
	"os"
	"encoding/json"
)

var musicNote MusicNote = MusicNote{}


func main() {
	// Print's the manual.
	fmt.Printf("%s", manual)
	time.Sleep(5 * time.Second)  // 5 seconds.

	// Debug: Only to test the JSON format.
	// str_json_test := MusicScoreToJsonString(music_01)
	// fmt.Printf("\n str_json_test: \n\n%s\n\n", str_json_test)

	jsonFilePathAndName := ""
	if len(os.Args) > 1{
		jsonFilePathAndName = os.Args[1]
	}

	// Reads the music from JSON file and writes the name to the screen.
	music_01 = getReadMusicScoreFromJSON(jsonFilePathAndName)
	fmt.Printf("\n\n\nMusic name: %s\n\n Description: %s\n", music_01.Name, music_01.Description )
	music_01.MSResetToRepeat()
	time.Sleep(2 * time.Second)  // 2 seconds.

	portaudio.Initialize()
	defer portaudio.Terminate()
	e := newMicophone(time.Second / 3)
	defer e.Close()
	chk(e.Start())
	time.Sleep(15* 60 * time.Second)  // 15 minuts ou Ctrl + C
	chk(e.Stop())
	//fmt.Printf("len %d, cap %d\n", input_buffer_len, input_buffer_cap)
}

type microphone struct {
	*portaudio.Stream
	buffer []float32
	i      int
}

func newMicophone(delay time.Duration) *microphone {
	// Inicializes the flute music notes.
	musicNote.MNnew()
	// Expand the music into a 2D array of runes with the music.
	music_01.MSExpandIntoArray()
	// Screen buffer where all text is written, before display.

	h, err := portaudio.DefaultHostApi()
	chk(err)
	p := portaudio.LowLatencyParameters(h.DefaultInputDevice, h.DefaultOutputDevice)
	p.Input.Channels = 1
	p.Output.Channels = 1
	e := &microphone{buffer: make([]float32, int(p.SampleRate*delay.Seconds()))}
	e.Stream, err = portaudio.OpenStream(p, e.processAudio)
	chk(err)
	return e
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

//#####################################
//#####################################

var input_buffer_len int = 0
var input_buffer_cap int = 0


var inputBuffer [buffLen]float32 = [buffLen]float32{}
var inputBufferIndex int = 0
var flag_counter int = 0

func (e *microphone) processAudio(in, out []float32) {

	// Samples rate is 44100.
    // fmt.Printf( "Sample rate:  %d\n", s.SampleRate())

	input_buffer_len = len(in)
	input_buffer_cap = cap(in)

	for i := range out {
		inputBuffer[inputBufferIndex] = in[i]

		if inputBufferIndex == buffLen - 1 {
			// Process one buffer then doesn't process the next buffer of length buffLen.
			if flag_counter == 0 {
				// Process buffer with algoritm.
				frequency, /*probability*/ _ := findMainFrequency(&inputBuffer)
				//fmt.Printf("Main Frequency: %f - Probability: %f \n", frequency, probability)
				// musicNote.MNPrintNote(frequency)

				// Writes the flute drawing in text into the screenBuffer
				playedNote := musicNote.MNPrintNoteToScreenBuffer(frequency)

				// Writes the sheet music into the screen.
				music_01.MSPrintMusicSheetToScreenBuffer(playedNote)
				// Makes the score move from the right to the left,
				music_01.MSUpdateMovement()

				// Writes screenBuffer to the screen with Printf.
				printScreenBuffer()

				flag_counter++
			} else {
				// Jumps each two buffers.
				if flag_counter < 2 {
					flag_counter++
				} else {
					flag_counter = 0
				}
			}
			inputBufferIndex = 0
		}else{
			inputBufferIndex++
		}
	}

}


///////////////////////////////////////////////////

const buffLen int = 5000 // 11025

func findMainFrequency(buff *[buffLen]float32 ) (frequency float64, probability float64){
	//arr :=  [YIN_SAMPLING_RATE / 2]float64{}
	//yin := Yin{0,0, arr, 0.0, 0.0}
	yin := Yin{}
	//bufferSize := 44100
	bufferSize := buffLen  // 11025
	threashold := 0.05
	yin.YinInit( bufferSize, threashold )
	frequency = yin.YinGetPitch(buff)
	probability = yin.YinGetProbability()
	return frequency, probability
}


//################################################################

const YIN_SAMPLING_RATE int = 44100
const YIN_DEFAULT_THRESHOLD float64 = 0.15
const BUFF_SIZE int = buffLen // 5000 // 11025 // 22050

type Yin struct {
	bufferSize     int     // Size of the buffer to process.
	halfBufferSize int     // Half of buffer size.
	yinBuffer      [BUFF_SIZE / 2]float64 // Buffer that stores the results of the intermediate processing steps of the algorithm
	probability    float64 // Probability that the pitch found is correct as a decimal (i.e 0.85 is 85%)
	threshold      float64 // Allowed uncertainty in the result as a decimal (i.e 0.15 is 15%)
}


// threshold  Allowed uncertainty (e.g 0.05 will return a pitch with ~95% probability)
func (Y *Yin) YinInit(bufferSize int, threshold float64) {
	// Initialise the fields of the Yin structure passed in.
	Y.bufferSize = bufferSize;
	Y.halfBufferSize = bufferSize / 2;
	Y.probability = 0.0;
	Y.threshold = threshold;

	// Allocate the autocorellation buffer and initialise it to zero.
	Y.yinBuffer = [BUFF_SIZE / 2]float64{}
}


// Runs the Yin pitch detection algortihm
//        buffer       - Buffer of samples to analyse
// return pitchInHertz - Fundamental frequency of the signal in Hz. Returns -1 if pitch can't be found
func (Y *Yin) YinGetPitch(buffer *[BUFF_SIZE]float32) (pitchInHertz float64) {
	//tauEstimate int      := -1
	pitchInHertz = -1

	// Step 1: CalcuYinGetPitchlates the squared difference of the signal with a shifted version of itself.
	Y.yinDifference(buffer)

	// Step 2: Calculate the cumulative mean on the normalised difference calculated in step 1.
	Y.yinCumulativeMeanNormalizedDifference()

	// Step 3: Search through the normalised cumulative mean array and find values that are over the threshold.
	tauEstimate := Y.yinAbsoluteThreshold()

	// Step 5: Interpolate the shift value (tau) to improve the pitch estimate.
	if(tauEstimate != -1){
		pitchInHertz = float64(YIN_SAMPLING_RATE) / Y.yinParabolicInterpolation(tauEstimate)
	}

	return pitchInHertz;
}


// Certainty of the pitch found
// return ptobability - Returns the certainty of the note found as a decimal (i.e 0.3 is 30%)
func (Y *Yin) YinGetProbability() (probability float64) {
	return Y.probability
}


// Step 1: Calculates the squared difference of the signal with a shifted version of itself.
// @param buffer Buffer of samples to process.
//
// This is the Yin algorithms tweak on autocorellation. Read http://audition.ens.fr/adc/pdf/2002_JASA_YIN.pdf
// for more details on what is in here and why it's done this way.
func (Y *Yin) yinDifference(buffer *[BUFF_SIZE]float32) {
	// Calculate the difference for difference shift values (tau) for the half of the samples.
	for tau := 0; tau < Y.halfBufferSize; tau++ {

		// Take the difference of the signal with a shifted version of itself, then square it.
		// (This is the Yin algorithm's tweak on autocorellation)
		for i := 0; i < Y.halfBufferSize; i++{
			delta := float64(buffer[i]) - float64(buffer[i + tau])
			Y.yinBuffer[tau] += delta * delta;
		}
	}
}


// Step 2: Calculate the cumulative mean on the normalised difference calculated in step 1
//
// This goes through the Yin autocorellation values and finds out roughly where shift is which
// produced the smallest difference
func (Y *Yin) yinCumulativeMeanNormalizedDifference() {
	runningSum := 0.0;
	Y.yinBuffer[0] = 1;

	// Sum all the values in the autocorellation buffer and nomalise the result, replacing
	// the value in the autocorellation buffer with a cumulative mean of the normalised difference.
	for tau := 1; tau < Y.halfBufferSize; tau++ {
		runningSum += Y.yinBuffer[tau]
		Y.yinBuffer[tau] *= float64(tau) / runningSum
	}
}


// Step 3: Search through the normalised cumulative mean array and find values that are over the threshold
// return Shift (tau) which caused the best approximate autocorellation. -1 if no suitable value is found
// over the threshold.
func (Y *Yin) yinAbsoluteThreshold() int {

	var tau int

	// Search through the array of cumulative mean values, and look for ones that are over the threshold
	// The first two positions in yinBuffer are always so start at the third (index 2)
	for tau = 2; tau < Y.halfBufferSize; tau++ {
		if (Y.yinBuffer[tau] < Y.threshold) {
			for (tau + 1 < Y.halfBufferSize) && (Y.yinBuffer[tau + 1] < Y.yinBuffer[tau]) {
				tau++;
			}

			/* found tau, exit loop and return
			 * store the probability
			 * From the YIN paper: The yin->threshold determines the list of
			 * candidates admitted to the set, and can be interpreted as the
			 * proportion of aperiodic power tolerated
			 * within a periodic signal.
			 *
			 * Since we want the periodicity and and not aperiodicity:
			 * periodicity = 1 - aperiodicity */
			Y.probability = 1 - Y.yinBuffer[tau];
			break;
		}
	}

	// if no pitch found, tau => -1
	if (tau == Y.halfBufferSize || Y.yinBuffer[tau] >= Y.threshold) {
		tau = -1;
		Y.probability = 0;
	}

	return tau;
}

// Step 5: Interpolate the shift value (tau) to improve the pitch estimate.
// tauEstimate [description]
// Return
// The 'best' shift value for autocorellation is most likely not an interger shift of the signal.
// As we only autocorellated using integer shifts we should check that there isn't a better fractional
// shift value.
func (Y *Yin) yinParabolicInterpolation(tauEstimate int) float64 {

	var betterTau float64
	var x0 int
	var x2 int

	// Calculate the first polynomial coeffcient based on the current estimate of tau.
	if tauEstimate < 1 {
		x0 = tauEstimate;
	} else {
		x0 = tauEstimate - 1;
	}

	// Calculate the second polynomial coeffcient based on the current estimate of tau.
	if tauEstimate + 1 < Y.halfBufferSize {
		x2 = tauEstimate + 1;
	} else {
		x2 = tauEstimate;
	}

	// Algorithm to parabolically interpolate the shift value tau to find a better estimate.
	if x0 == tauEstimate {
		if Y.yinBuffer[tauEstimate] <= Y.yinBuffer[x2] {
			betterTau = float64(tauEstimate);
		} else {
			betterTau = float64(x2);
		}
	} else if x2 == tauEstimate {
		if Y.yinBuffer[tauEstimate] <= Y.yinBuffer[x0]{
			betterTau = float64(tauEstimate);
		} else {
			betterTau = float64(x0);
		}
	} else {
		var s0, s1, s2 float64
		s0 = Y.yinBuffer[x0];
		s1 = Y.yinBuffer[tauEstimate];
		s2 = Y.yinBuffer[x2];
		// fixed AUBIO implementation, thanks to Karl Helgason:
		// (2.0f * s1 - s2 - s0) was incorrectly multiplied with -1
		betterTau = float64(tauEstimate) + (s2 - s0) / (2 * (2 * s1 - s2 - s0));
	}

	return betterTau;
}


const (
	EMPTY int = iota
	DO
	RE
	MI
	FA
	SOL
	LA
	SI
	DO_HIGH
)


const fluteNoteLen int = 9

type MusicNote struct {
	//description     [fluteNoteLen]string     // Description
	note            [fluteNoteLen]string       // Name of the music note.
	frequency       [fluteNoteLen]int          // Frequency of the music note.
	textFluteOutput [fluteNoteLen][13]string   // Text representation of the flute drawing.
    VisualIndex     [fluteNoteLen]int          // The index that shows visualy in the Music Score for this note.
}

func (MN *MusicNote) MNnew() {

	noteIndex := 0
	// Recorder with no holes covered.
	MN.note[noteIndex]      = "Recorder with no holes covered."
	MN.frequency[noteIndex] =  1180
	MN.VisualIndex[noteIndex] = 0
	MN.textFluteOutput[noteIndex][ 0] = "  ---    "
	MN.textFluteOutput[noteIndex][ 1] = " | = |   "
	MN.textFluteOutput[noteIndex][ 2] = " |   |   "
	MN.textFluteOutput[noteIndex][ 3] = " | O | O "
	MN.textFluteOutput[noteIndex][ 4] = " | O |   "
	MN.textFluteOutput[noteIndex][ 5] = " | O |   "
	MN.textFluteOutput[noteIndex][ 6] = " | O |   "
	MN.textFluteOutput[noteIndex][ 7] = " | O |   "
	MN.textFluteOutput[noteIndex][ 8] = " | O |   "
	MN.textFluteOutput[noteIndex][ 9] = " |O  |   "
	MN.textFluteOutput[noteIndex][10] = "  | |    "
	MN.textFluteOutput[noteIndex][11] = " |   |   "
	MN.textFluteOutput[noteIndex][12] = "  ---    "


	noteIndex = 1
	// Do - All holes covered.
	MN.note[noteIndex]      = "Do"
	MN.frequency[noteIndex] =  521  // The correct should be 527 Hz, but isn't this one that comes out.
	MN.VisualIndex[noteIndex] = 9
	MN.textFluteOutput[noteIndex][ 0] = "  ---    "
	MN.textFluteOutput[noteIndex][ 1] = " | = |   "
	MN.textFluteOutput[noteIndex][ 2] = " |   |   "
	MN.textFluteOutput[noteIndex][ 3] = " | # | # "
	MN.textFluteOutput[noteIndex][ 4] = " | # |   "
	MN.textFluteOutput[noteIndex][ 5] = " | # |   "
	MN.textFluteOutput[noteIndex][ 6] = " | # |   "
	MN.textFluteOutput[noteIndex][ 7] = " | # |   "
	MN.textFluteOutput[noteIndex][ 8] = " | # |   "
	MN.textFluteOutput[noteIndex][ 9] = " |#  |   "
	MN.textFluteOutput[noteIndex][10] = "  | |    "
	MN.textFluteOutput[noteIndex][11] = " |   |   "
	MN.textFluteOutput[noteIndex][12] = "  ---    "


	noteIndex = 2
	// Re - All holes covered except the last one.
	MN.note[noteIndex]      = "Re ---A1#/B1b"
	MN.frequency[noteIndex] =  630
	MN.VisualIndex[noteIndex] = 8
	MN.textFluteOutput[noteIndex][ 0] = "  ---    "
	MN.textFluteOutput[noteIndex][ 1] = " | = |   "
	MN.textFluteOutput[noteIndex][ 2] = " |   |   "
	MN.textFluteOutput[noteIndex][ 3] = " | # | # "
	MN.textFluteOutput[noteIndex][ 4] = " | # |   "
	MN.textFluteOutput[noteIndex][ 5] = " | # |   "
	MN.textFluteOutput[noteIndex][ 6] = " | # |   "
	MN.textFluteOutput[noteIndex][ 7] = " | # |   "
	MN.textFluteOutput[noteIndex][ 8] = " | # |   "
	MN.textFluteOutput[noteIndex][ 9] = " |O  |   "
	MN.textFluteOutput[noteIndex][10] = "  | |    "
	MN.textFluteOutput[noteIndex][11] = " |   |   "
	MN.textFluteOutput[noteIndex][12] = "  ---    "


	noteIndex = 3
	// Mi - All holes covered except the last two.
	MN.note[noteIndex]      = "Mi --- A1"
	MN.frequency[noteIndex] =  652  // The correct should be 627 Hz.
	MN.VisualIndex[noteIndex] = 7
	MN.textFluteOutput[noteIndex][ 0] = "  ---    "
	MN.textFluteOutput[noteIndex][ 1] = " | = |   "
	MN.textFluteOutput[noteIndex][ 2] = " |   |   "
	MN.textFluteOutput[noteIndex][ 3] = " | # | # "
	MN.textFluteOutput[noteIndex][ 4] = " | # |   "
	MN.textFluteOutput[noteIndex][ 5] = " | # |   "
	MN.textFluteOutput[noteIndex][ 6] = " | # |   "
	MN.textFluteOutput[noteIndex][ 7] = " | # |   "
	MN.textFluteOutput[noteIndex][ 8] = " | O |   "
	MN.textFluteOutput[noteIndex][ 9] = " |O  |   "
	MN.textFluteOutput[noteIndex][10] = "  | |    "
	MN.textFluteOutput[noteIndex][11] = " |   |   "
	MN.textFluteOutput[noteIndex][12] = "  ---    "


	noteIndex = 4
	// Fá - All holes covered except the last three.
	MN.note[noteIndex]      = "Fá  --- G1"
	MN.frequency[noteIndex] =  700 // Should be the frequency of 704 Hz.
	MN.VisualIndex[noteIndex] = 6
	MN.textFluteOutput[noteIndex][ 0] = "  ---    "
	MN.textFluteOutput[noteIndex][ 1] = " | = |   "
	MN.textFluteOutput[noteIndex][ 2] = " |   |   "
	MN.textFluteOutput[noteIndex][ 3] = " | # | # "
	MN.textFluteOutput[noteIndex][ 4] = " | # |   "
	MN.textFluteOutput[noteIndex][ 5] = " | # |   "
	MN.textFluteOutput[noteIndex][ 6] = " | # |   "
	MN.textFluteOutput[noteIndex][ 7] = " | O |   "
	MN.textFluteOutput[noteIndex][ 8] = " | O |   "
	MN.textFluteOutput[noteIndex][ 9] = " |O  |   "
	MN.textFluteOutput[noteIndex][10] = "  | |    "
	MN.textFluteOutput[noteIndex][11] = " |   |   "
	MN.textFluteOutput[noteIndex][12] = "  ---    "


	noteIndex = 5
	// Sol - With the first three holes covered.
	MN.note[noteIndex]      = "Sol --- F1"
	MN.frequency[noteIndex] =  780  // Should be 790 Hz.
	MN.VisualIndex[noteIndex] = 5
	MN.textFluteOutput[noteIndex][ 0] = "  ---    "
	MN.textFluteOutput[noteIndex][ 1] = " | = |   "
	MN.textFluteOutput[noteIndex][ 2] = " |   |   "
	MN.textFluteOutput[noteIndex][ 3] = " | # | # "
	MN.textFluteOutput[noteIndex][ 4] = " | # |   "
	MN.textFluteOutput[noteIndex][ 5] = " | # |   "
	MN.textFluteOutput[noteIndex][ 6] = " | O |   "
	MN.textFluteOutput[noteIndex][ 7] = " | O |   "
	MN.textFluteOutput[noteIndex][ 8] = " | O |   "
	MN.textFluteOutput[noteIndex][ 9] = " |O  |   "
	MN.textFluteOutput[noteIndex][10] = "  | |    "
	MN.textFluteOutput[noteIndex][11] = " |   |   "
	MN.textFluteOutput[noteIndex][12] = "  ---    "


	noteIndex = 6
	// La - With the first two holes covered.
	MN.note[noteIndex]      = "La --- E1"
	MN.frequency[noteIndex] =  882   // Should be 837 Hz.
	MN.VisualIndex[noteIndex] = 4
	MN.textFluteOutput[noteIndex][ 0] = "  ---    "
	MN.textFluteOutput[noteIndex][ 1] = " | = |   "
	MN.textFluteOutput[noteIndex][ 2] = " |   |   "
	MN.textFluteOutput[noteIndex][ 3] = " | # | # "
	MN.textFluteOutput[noteIndex][ 4] = " | # |   "
	MN.textFluteOutput[noteIndex][ 5] = " | O |   "
	MN.textFluteOutput[noteIndex][ 6] = " | O |   "
	MN.textFluteOutput[noteIndex][ 7] = " | O |   "
	MN.textFluteOutput[noteIndex][ 8] = " | O |   "
	MN.textFluteOutput[noteIndex][ 9] = " |O  |   "
	MN.textFluteOutput[noteIndex][10] = "  | |    "
	MN.textFluteOutput[noteIndex][11] = " |   |   "
	MN.textFluteOutput[noteIndex][12] = "  ---    "


	noteIndex = 7
	// Si - With the first hole covered.
	MN.note[noteIndex]      = "Si ---- D1"
	MN.frequency[noteIndex] =  985   // Shoud be 939 Hz.
	MN.VisualIndex[noteIndex] = 3
	MN.textFluteOutput[noteIndex][ 0] = "  ---    "
	MN.textFluteOutput[noteIndex][ 1] = " | = |   "
	MN.textFluteOutput[noteIndex][ 2] = " |   |   "
	MN.textFluteOutput[noteIndex][ 3] = " | # | # "
	MN.textFluteOutput[noteIndex][ 4] = " | O |   "
	MN.textFluteOutput[noteIndex][ 5] = " | O |   "
	MN.textFluteOutput[noteIndex][ 6] = " | O |   "
	MN.textFluteOutput[noteIndex][ 7] = " | O |   "
	MN.textFluteOutput[noteIndex][ 8] = " | O |   "
	MN.textFluteOutput[noteIndex][ 9] = " |O  |   "
	MN.textFluteOutput[noteIndex][10] = "  | |    "
	MN.textFluteOutput[noteIndex][11] = " |   |   "
	MN.textFluteOutput[noteIndex][12] = "  ---    "


	noteIndex = 8
	// Do high - With only two holes covered.
	MN.note[noteIndex]      = "Do high"
	MN.frequency[noteIndex] =  1040   // It should be 1054 Hz.
	MN.VisualIndex[noteIndex] = 4
	MN.textFluteOutput[noteIndex][ 0] = "  ---    "
	MN.textFluteOutput[noteIndex][ 1] = " | = |   "
	MN.textFluteOutput[noteIndex][ 2] = " |   |   "
	MN.textFluteOutput[noteIndex][ 3] = " | O | # "
	MN.textFluteOutput[noteIndex][ 4] = " | # |   "
	MN.textFluteOutput[noteIndex][ 5] = " | O |   "
	MN.textFluteOutput[noteIndex][ 6] = " | O |   "
	MN.textFluteOutput[noteIndex][ 7] = " | O |   "
	MN.textFluteOutput[noteIndex][ 8] = " | O |   "
	MN.textFluteOutput[noteIndex][ 9] = " |O  |   "
	MN.textFluteOutput[noteIndex][10] = "  | |    "
	MN.textFluteOutput[noteIndex][11] = " |   |   "
	MN.textFluteOutput[noteIndex][12] = "  ---    "

}

func (MN *MusicNote) MNPrintNote(frequency float64) {
	bestIndex := MN.MNFindFluteNoteIndex(frequency)

	var buffer bytes.Buffer
	for i:=0; i<13; i++ {
		buffer.WriteString(MN.textFluteOutput[bestIndex][i])
		buffer.WriteString("\n")
	}

	str_flute := buffer.String()

	//fmt.Printf("%d --- %d\n", bestIndex, MN.frequency[bestIndex])
	fmt.Printf("\n\n\n  %s\n%s", MN.note[bestIndex], str_flute)
}

func (MN *MusicNote) MNFindFluteNoteIndex(frequency float64) (bestIndex int) {
	bestIndex = -1
	lowestDelta := 9999999999.0

	// If frequency is -1.0 then no frequency was detected.
	if frequency + 1 > 0.0001 {
		for i := 0; i < fluteNoteLen; i++ {
			difference := frequency - float64(MN.frequency[i])
			delta := math.Sqrt(difference * difference)
			if delta < lowestDelta {
				lowestDelta = delta
				bestIndex = i
			}
		}
	}else{
		bestIndex = 0
	}

	return bestIndex
}

func (MN *MusicNote) MNPrintNoteToScreenBuffer(frequency float64) (playedNote int) {
	bestIndex := MN.MNFindFluteNoteIndex(frequency)

	for i:=0; i<13; i++ {
		runesFluteLine := []rune(MN.textFluteOutput[bestIndex][i])
		for j, e := range runesFluteLine{
			screenBuffer[i][j] = e
		}
	}
	return bestIndex
}

////////////////////////////////////////////
////////////////////////////////////////////
////////////////////////////////////////////

/*
// note enum:
const (
	VAZIO int = iota
	DO
	RE
	MI
	FA
	SOL
	LA
	SI
	DO_AGUDO
)
*/

type PlayNote struct {
	Note     int   `json:"note"`      // The note that's going to be played.
	Duration int   `json:"duration"`  // Number os steps it's vallid.
}

type MusicScore struct {
	Name               string      `json:"name"`        // Music score name.
	NotesList          []PlayNote  `json:"notesList"`   // Musical notes.
	Description        string      `json:"description"`
	duration           int
	expandedRunesArray [][]rune      // Expanded array of runes for the sheet music.
	indexSourceStart   int           // Index on the expandedRunesArray of the Start position. Copies from this position on the expandedRunesArray to the screenBuffer.
	indexTargetStart   int           // Index on the screenBuffer of the Start position.
	}


func (MS * MusicScore) MSDuration() (duration int) {
	duration = 0
	for _, e := range MS.NotesList {
		duration += e.Duration
	}
	return duration
}


func (MS * MusicScore) MSExpandIntoArray() {

	duration := MS.MSDuration()
	MS.duration = duration

	// Allocates memory for the 2D array [13](duration]
	//var MSTextArray [][]rune = [][]rune{}
	MSTextArray := make([][]rune, 13)
	for i := range MSTextArray {
		MSTextArray[i] = make([]rune, duration)
	}

	// Initialize buffer.
	for i:=3; i<13-3; i++ {
		for j:=0; j<duration; j++ {
			MSTextArray[i][j] = '.'
		}
	}

	currentPos := 0
	for _, e := range MS.NotesList{
		switch e.Note {
		case 	EMPTY:

		/*
        case  DO:
		case  RE:
		case  MI:
		case  FA:
		case  SOL:
		case  LA:
		case  SI:
		*/
		case  DO_HIGH:
			index := musicNote.VisualIndex[e.Note]
			MSTextArray[index][currentPos] = 'D'
			for i:=0; i<e.Duration-1; i++{
				currentPos++
				MSTextArray[index][currentPos] = '_'
			}

		default:
			// Processes the normal notes!
			index := musicNote.VisualIndex[e.Note]
			MSTextArray[index][currentPos] = 'S'
			for i:=0; i<e.Duration-1; i++{
				currentPos++
				MSTextArray[index][currentPos] = '_'
			}
		}

		currentPos++
	}

	MS.expandedRunesArray = MSTextArray
}

func (MS *MusicScore) MSPrintMusicSheetToScreenBuffer(note int) {

	// Initialize the screen with '.'
	for i:=3; i<NUM_LINES_SCREEN - 3; i++ {
		for j:=10; j<MAX_SCREEN_WIDE; j++ {
			screenBuffer[i][j] = '.'
		}
	}

	currentIndexSourceStart := MS.indexSourceStart
	for indexTarget:=MS.indexTargetStart; indexTarget < MAX_SCREEN_WIDE; indexTarget++{
		for i:=3; i<NUM_LINES_SCREEN - 3; i++{
			if currentIndexSourceStart < MS.duration {
				screenBuffer[i][indexTarget] = MS.expandedRunesArray[i][currentIndexSourceStart]
			}
		}
		currentIndexSourceStart++
	}

//	indexSourceStart   int           // Index on the expandedRunesArray of the Start position. Copies from this position on the expandedRunesArray to the screenBuffer.
//	indexTargetStart   int           // Index on the screenBuffer of the Start position.

	visualIndex := musicNote.VisualIndex[note]

	//fmt.Printf("visualIndex: %d", visualIndex)

	// Draw the vertical line (Win Line) on the left of the screen that markes where the notes are scorred.
	for i:=3; i<NUM_LINES_SCREEN - 3; i++ {

		// It validates if the played musical note was correct or if it was an incorrect note.
		// Knowing that it marks it with the simble 'X', '@' or '|'.

		rune := screenBuffer[i][10]

		if visualIndex == i && (rune == 'S' || rune == '_' || rune == 'D') {
			//fmt.Printf("rune: %s", string(rune) )

			//if rune == 'S' || rune == '_' || rune == 'D'{
				if note == DO_HIGH && rune != 'D' {
					screenBuffer[i][10] = '@'
					if currentScore > 0 {
						// Each wrong note decreases the score by one.
						currentScore--
					}
					continue
				}
				screenBuffer[i][10] = 'X'
				// Each right note increases the score by ten.
				currentScore += 10
			//}

		}else{
			underRune := screenBuffer[i][10]
			if underRune == 'S' || underRune == '_' || underRune == 'D'{
				screenBuffer[i][10] = '@'
				if currentScore > 0 {
					// Each wrong note decreases the score by one.
					currentScore--
				}
			}else {
				screenBuffer[i][10] = '|'
			}
		}
	}

}

func (MS *MusicScore) MSUpdateMovement() {

	if MS.indexTargetStart > 10 {
		MS.indexTargetStart--
	}
	if MS.indexTargetStart == 10{
	   if MS.duration > MS.indexSourceStart {
		   MS.indexSourceStart++
	   }else{
		   MS.MSResetToRepeat()
	   }
	}

	// indexSourceStart   int           // Index on the expandedRunesArray of the Start position. Copies from this position on the expandedRunesArray to the screenBuffer.
	// indexTargetStart   int           // Index on the screenBuffer of the Start position.

	// Debug:
	// fmt.Printf("MS.indexSourceStart: %d MS.indexTargetStart: %d\n", MS.indexSourceStart, MS.indexTargetStart)
}

// Reset's the state and repeats the music score.
func (MS *MusicScore) MSResetToRepeat() {
	MS.indexSourceStart = 0
	MS.indexTargetStart = MAX_SCREEN_WIDE
}


///////////////////////////////////////////////
///////////////////////////////////////////////
///////////////////////////////////////////////


const NUM_LINES_SCREEN int = 13
const MAX_SCREEN_WIDE int  = 40

// Screen buffer where all text is written, before display.
var screenBuffer [NUM_LINES_SCREEN][MAX_SCREEN_WIDE]rune = [NUM_LINES_SCREEN][MAX_SCREEN_WIDE]rune{}
var currentScore int = 0

func printScreenBuffer(){

	// Joins every rune in a array.
	runeArray := []rune{}
	for i:=0; i<NUM_LINES_SCREEN; i++ {
		for j:=0; j<MAX_SCREEN_WIDE; j++ {
			runeArray = append(runeArray, screenBuffer[i][j])
		}
		runeArray = append(runeArray, '\n')
	}

	myStr := string(runeArray)
	fmt.Printf("\n\n\n\n                Galileu's Flute\n\n                    Score: %d\n%s", currentScore, myStr)
}

func getReadMusicScoreFromJSON(jsonFilePathAndName string) MusicScore {
	// If the path comes empty fill it with the default file path name.
	if jsonFilePathAndName == "" {
		jsonFilePathAndName = "./music_01.json"
	}
	raw, err := ioutil.ReadFile(jsonFilePathAndName)
	if err != nil {
		fmt.Println("Error Reading the JSON Music Score file!")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var musicScore MusicScore
	err = json.Unmarshal(raw, &musicScore)
	if err != nil {
		fmt.Println("Error ahile parsing the JSON Music Score file!")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return musicScore
}

func MusicScoreToJsonString(p interface{}) string {
	bytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return string(bytes)
}


var manual string = `

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


Example of the output:


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

`


//##################
//# Music Score 01
//#

var music_01 MusicScore = MusicScore {
	Name : "Music 01 !",
	NotesList :
	[]PlayNote{
		{Note: SI,
			Duration: 1},

		{Note: EMPTY,
			Duration: 1},

		{Note: SI,
			Duration: 1},

		{Note: EMPTY,
			Duration: 2},

		{Note: SI,
			Duration: 5},

		{Note: LA,
			Duration: 2},

		{Note: EMPTY,
			Duration: 1},

		{Note: SOL,
			Duration: 2},

		{Note: FA,
			Duration: 2},

		{Note: MI,
			Duration: 2},

		{Note: RE,
			Duration: 2},

		{Note: DO,
			Duration: 2},

		{Note: FA,
			Duration: 2},

		{Note: LA,
			Duration: 2},

		{Note: DO_HIGH,
			Duration: 1},

	},
	Description: "This is a music invented by me without much skill in music :-)",
	indexSourceStart : 0,                       // Index on the expandedRunesArray of the Start position. Copies from this position on the expandedRunesArray to the screenBuffer.
	indexTargetStart : MAX_SCREEN_WIDE - 1 ,    // Index on the screenBuffer of the Start position.
}



