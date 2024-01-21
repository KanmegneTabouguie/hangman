// src/Game.js
import React, { useState, useEffect } from 'react';
import 'bootstrap/dist/css/bootstrap.min.css';

const Game = () => {
  const [word, setWord] = useState('');
  const [guessedLetters, setGuessedLetters] = useState([]);
  const [lives, setLives] = useState(0);

  useEffect(() => {
    // Fetch initial game state when the component mounts
    fetch('http://localhost:8082/init?playerName=Player1&difficulty=medium')
      .then(response => response.json())
      .then(data => {
        setWord(data.word);
        setGuessedLetters(data.guessedLetters);
        setLives(data.lives);
      })
      .catch(error => console.error('Error fetching initial game state:', error));
  }, []);

  const handleGuess = guess => {
    // Send guess to the server and update the game state
    fetch('http://localhost:8082/guess', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ guess }),
    })
      .then(response => response.json())
      .then(data => {
        setWord(data.word);
        setGuessedLetters(data.guessedLetters);
        setLives(data.lives);
      })
      .catch(error => console.error('Error making a guess:', error));
  };

  return (
    <div className="container mt-4">
      <h1>Hangman Game</h1>
      <div className="mb-4">
        <p>Word: {word}</p>
        <p>Guessed Letters: {guessedLetters.join(', ')}</p>
        <p>Lives: {lives}</p>
      </div>
      <div>
        <label htmlFor="guessInput">Enter a letter: </label>
        <input
          type="text"
          id="guessInput"
          maxLength="1"
          onChange={e => handleGuess(e.target.value)}
        />
      </div>
    </div>
  );
};

export default Game;
