/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

/*
Sprite sheet files start with a 2 byte integer telling the number of frames inside the file,
followed by the sprite dimensions as 1 byte width and height. Next is a list of all frames,
starting with their y and x offset, followed by width and height, each as 1 byte value.
Last comes the offset of the frame inside the file, stored as 4 byte integer.
If the width times height is greater than the difference between this and the next
offset, then the frame is compressed as specified below. Else it is to be read as a usual
indexed 256 color bitmap.
*/

package resource
