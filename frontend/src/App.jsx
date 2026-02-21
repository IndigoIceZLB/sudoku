import { useState, useEffect } from 'react';
import axios from 'axios';
import './App.css';

// 你的 Render 后端地址 (注意：末尾不要带斜杠)
const API_URL = "https://sudokuapi-rlim.onrender.com";

function App() {
  const [board, setBoard] = useState([]); // 当前棋盘状态
  const [initialBoard, setInitialBoard] = useState([]); // 初始题目（用于判断哪些格子是锁定的）
  const [loading, setLoading] = useState(false);
  const [difficulty, setDifficulty] = useState('easy');

  // 获取新游戏
  const fetchNewGame = async (level) => {
    setLoading(true);
    try {
      const res = await axios.get(`${API_URL}/api/new-game?level=${level}`);
      // 后端返回的 0 代表空格，前端为了好看可以转成空字符串，或者保持 0
      setBoard(res.data.puzzle);
      setInitialBoard(JSON.parse(JSON.stringify(res.data.puzzle))); // 深拷贝一份作为初始状态
    } catch (error) {
      console.error("Failed to fetch game:", error);
      alert("无法连接到服务器，请检查网络或稍后再试。");
    }
    setLoading(false);
  };

  // 页面加载时自动开始一局
  useEffect(() => {
    fetchNewGame('easy');
  }, []);

  // 处理输入
  const handleInputChange = (rowIndex, colIndex, value) => {
    // 如果是初始题目中的数字，不允许修改
    if (initialBoard[rowIndex][colIndex] !== 0) return;

    // 只允许输入 1-9 的数字
    const num = parseInt(value);
    if (value === '' || (num >= 1 && num <= 9)) {
      const newBoard = [...board];
      newBoard[rowIndex] = [...newBoard[rowIndex]]; // 浅拷贝这一行
      newBoard[rowIndex][colIndex] = value === '' ? 0 : num;
      setBoard(newBoard);
    }
  };

  return (
    <div className="container">
      <h1>Sudoku Go</h1>
      
      <div className="controls">
        <select value={difficulty} onChange={(e) => setDifficulty(e.target.value)}>
          <option value="easy">Easy</option>
          <option value="medium">Medium</option>
          <option value="hard">Hard</option>
        </select>
        <button onClick={() => fetchNewGame(difficulty)} disabled={loading}>
          {loading ? "Loading..." : "New Game"}
        </button>
      </div>

      <div className="board">
        {board.map((row, rowIndex) => (
          <div key={rowIndex} className="row">
            {row.map((cell, colIndex) => {
              const isInitial = initialBoard[rowIndex][colIndex] !== 0;
              return (
                <input
                  key={`${rowIndex}-${colIndex}`}
                  type="text"
                  maxLength="1"
                  className={`cell ${isInitial ? 'initial' : ''}`}
                  value={cell === 0 ? '' : cell}
                  readOnly={isInitial}
                  onChange={(e) => handleInputChange(rowIndex, colIndex, e.target.value)}
                />
              );
            })}
          </div>
        ))}
      </div>
    </div>
  );
}

export default App;