export const WORD_LENGTH = 5;
export const MAX_ATTEMPTS = 6;

export type Color = 'grey' | 'green' | 'yellow';

export type Letter = {
  color: Color;
  letter: string;
};

export type Attempt = {
  word_guess: string;
  letters: [Letter, Letter, Letter, Letter, Letter];
};
