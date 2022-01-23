import './LetterGrid.css';

import * as React from 'react';

import { Attempt, MAX_ATTEMPTS, WORD_LENGTH } from '../types';

type LetterGridProps = {
  attempts: Attempt[];
  input: string;
};

export const LetterGrid = ({ attempts, input }: LetterGridProps) => {
  return (
    <div className="container">
      {Array(MAX_ATTEMPTS)
        .fill(0)
        .map((_, row) => {
          if (attempts[row]) {
            return attempts[row].letters.map((letter, i) => (
              <div className={`cell cell-${row}-${i} color-${letter.color}`} key={`${row}-${i}`}>
                {letter.letter}
              </div>
            ));
          }
          if (row === attempts.length) {
            return Array(WORD_LENGTH)
              .fill(0)
              .map((_, i) => (
                <div className={`cell cell-${row}-${i}`} key={`${row}-${i}`}>
                  {input[i]}
                </div>
              ));
          }
          return Array(WORD_LENGTH)
            .fill(0)
            .map((_, i) => <div className={`cell cell-${row}-${i}`} key={`${row}-${i}`}></div>);
        })}
    </div>
  );
};
