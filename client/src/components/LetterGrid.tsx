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
          if (row === attempts.length - 1) {
            return attempts[row].letters.map((letter, i) => (
              <div className={`cell-container column-${i}`} key={`${row}-${i}`}>
                <div className="cell-card">
                  <div className={`cell cell-card-back color-${letter.color}`}>{letter.letter}</div>
                  <div className="cell cell-card-front">{letter.letter}</div>
                </div>
              </div>
            ));
          }
          if (attempts[row]) {
            return attempts[row].letters.map((letter, i) => (
              <div className="cell-container" key={`${row}-${i}`}>
                <div className={`cell color-${letter.color}`}>{letter.letter}</div>
              </div>
            ));
          }
          if (row === attempts.length) {
            return Array(WORD_LENGTH)
              .fill(0)
              .map((_, i) => (
                <div className="cell-container" key={`${row}-${i}`}>
                  <div className={`${input[i] ? 'cell' : 'cell-empty'}`}>{input[i]}</div>
                </div>
              ));
          }
          return Array(WORD_LENGTH)
            .fill(0)
            .map((_, i) => (
              <div className="cell-container" key={`${row}-${i}`}>
                <div className="cell-empty"></div>
              </div>
            ));
        })}
    </div>
  );
};
