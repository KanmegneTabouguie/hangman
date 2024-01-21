// src/Game.js
import React, { useState, useEffect } from 'react';
import 'bootstrap/dist/css/bootstrap.min.css';

const Game = () => {
  const [word, setWord] = useState('');
  const [guessedLetters, setGuessedLetters] = useState([]);
  const [lives, setLives] = useState(0);
  const [gameStatus, setGameStatus] = useState('');
  const [isGameOver, setIsGameOver] = useState(false);

  useEffect(() => {
    // Fetch initial game state when the component mounts
    fetch('http://localhost:8083/init?playerName=Player1&difficulty=medium')
      .then(response => response.json())
      .then(data => {
        setWord(data.word);
        setGuessedLetters(data.guessedLetters);
        setLives(data.lives);
        setGameStatus('');
        setIsGameOver(false);
      })
      .catch(error => console.error('Error fetching initial game state:', error));
  }, []);

  const handleGuess = guess => {
    // Check if the game is already over
    if (isGameOver) {
      console.log('The game is already over. Please start a new game.');
      return;
    }

    // Validate the guess (single letter)
    if (!/^[a-zA-Z]$/.test(guess)) {
      console.log('Invalid guess. Please enter a single letter.');
      return;
    }

    // Send guess to the server and update the game state
    fetch('http://localhost:8083/guess', {
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

        // Check if the game is over
        if (data.word.indexOf('_') === -1) {
          setGameStatus('Congratulations! You guessed the word.');
          setIsGameOver(true);
        } else if (data.lives === 0) {
          setGameStatus('Game over! The word was: ' + data.word);
          setIsGameOver(true);
        } else {
          setGameStatus('');
        }
      })
      .catch(error => console.error('Error making a guess:', error));
  };

  const startNewGame = () => {
    // Fetch a new game state
    fetch('http://localhost:8083/init?playerName=Player1&difficulty=medium')
      .then(response => response.json())
      .then(data => {
        setWord(data.word);
        setGuessedLetters(data.guessedLetters);
        setLives(data.lives);
        setGameStatus('');
        setIsGameOver(false);
      })
      .catch(error => console.error('Error starting a new game:', error));
  };

  return (
    <div className="container mt-4">
      <h1>Hangman Game</h1>
      <div className="mb-4">
        <p>Word: {word}</p>
        <p>Guessed Letters: {guessedLetters.join(', ')}</p>
        <p>Lives: {lives}</p>
        <p>{gameStatus}</p>
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
      <div className="mt-3">
        <button className="btn btn-primary" onClick={startNewGame}>
          Start New Game
        </button>
      </div>
    </div>
  );
};

export default Game;
