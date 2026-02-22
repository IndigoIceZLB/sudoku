import { useState, useEffect, useRef } from 'react';
import axios from 'axios';
import './App.css';

// 你的 Render 后端地址
const API_URL = "https://sudokuapi-rlim.onrender.com";

function App() {
  const [board, setBoard] = useState([]); 
  const [initialBoard, setInitialBoard] = useState([]); 
  const [solution, setSolution] = useState([]);
  const [loading, setLoading] = useState(false);
  const [difficulty, setDifficulty] = useState('easy');
  
  const [timer, setTimer] = useState(0);
  const [isGameActive, setIsGameActive] = useState(false);
  const [isWon, setIsWon] = useState(false);
  const [username, setUsername] = useState('');
  const [leaderboard, setLeaderboard] = useState([]);
  const [showLeaderboard, setShowLeaderboard] = useState(false);

  const timerRef = useRef(null);

  const formatTime = (seconds) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
  };

  // 🛑 核心修复：停止计时器的辅助函数
  const stopTimer = () => {
    if (timerRef.current) {
      clearInterval(timerRef.current);
      timerRef.current = null;
    }
  };

  const fetchNewGame = async (level) => {
    // 1. 开始请求前，先把旧的定时器关掉！(修复双倍速问题)
    stopTimer();
    
    setLoading(true);
    setIsGameActive(false);
    setIsWon(false);
    setTimer(0);
    setShowLeaderboard(false);

    try {
      const res = await axios.get(`${API_URL}/api/new-game?level=${level}`);
      
      // 2. 再次确保没有残留定时器
      stopTimer();

      setBoard(res.data.puzzle);
      setInitialBoard(JSON.parse(JSON.stringify(res.data.puzzle)));
      setSolution(res.data.solution); 
      
      // 3. 启动新定时器
      setIsGameActive(true);
      timerRef.current = setInterval(() => {
        setTimer((prev) => prev + 1);
      }, 1000);

    } catch (error) {
      console.error("Failed to fetch game:", error);
      alert("无法连接到服务器");
    }
    setLoading(false);
  };

  const checkWin = (currentBoard) => {
    if (JSON.stringify(currentBoard) === JSON.stringify(solution)) {
      // 🛑 核心修复：胜利瞬间立刻停止计时
      stopTimer();
      setIsGameActive(false);
      setIsWon(true);
      fetchLeaderboard(difficulty);
    }
  };

  const handleInputChange = (rowIndex, colIndex, value) => {
    if (!isGameActive) return;
    if (initialBoard[rowIndex][colIndex] !== 0) return;

    const num = parseInt(value);
    if (value === '' || (num >= 1 && num <= 9)) {
      const newBoard = JSON.parse(JSON.stringify(board));
      newBoard[rowIndex][colIndex] = value === '' ? 0 : num;
      setBoard(newBoard);
      
      const hasEmpty = newBoard.some(row => row.includes(0));
      if (!hasEmpty) {
        checkWin(newBoard);
      }
    }
  };

  const submitScore = async () => {
    if (!username) return alert("请输入名字！");
    try {
      // 发送请求
      await axios.post(`${API_URL}/api/submit-score`, {
        username,
        difficulty,
        time_spent: timer // 注意这里用的是停止后的 timer 值
      });
      alert("分数提交成功！");
      setIsWon(false);
      fetchLeaderboard(difficulty);
      setShowLeaderboard(true);
    } catch (error) {
      // 打印详细错误到控制台，方便调试
      console.error("Submit Error:", error.response ? error.response.data : error.message);
      alert("提交失败，请按 F12 打开控制台(Console)查看具体错误原因");
    }
  };

  const fetchLeaderboard = async (diff) => {
    try {
      const res = await axios.get(`${API_URL}/api/leaderboard?difficulty=${diff}`);
      setLeaderboard(res.data.leaderboard || []);
    } catch (error) {
      console.error(error);
    }
  };

  useEffect(() => {
    fetchNewGame('easy');
    return () => stopTimer(); // 组件卸载时清理
  }, []);

  return (
    <div className="container">
      <h1>Sudoku Go</h1>
      
      <div className="header-info">
        <div className="timer">⏱️ {formatTime(timer)}</div>
        <button onClick={() => {
          fetchLeaderboard(difficulty);
          setShowLeaderboard(true);
        }}>🏆 排行榜</button>
      </div>

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

      {isWon && (
        <div className="modal-overlay">
          <div className="modal">
            <h2>🎉 You Won! 🎉</h2>
            <p>Difficulty: {difficulty}</p>
            {/* 显示最终定格的时间 */}
            <p>Time: {formatTime(timer)}</p>
            <input 
              type="text" 
              placeholder="Enter your name" 
              value={username}
              onChange={(e) => setUsername(e.target.value)}
            />
            <button onClick={submitScore}>Submit Score</button>
            <button onClick={() => setIsWon(false)} className="close-btn">Close</button>
          </div>
        </div>
      )}

      {showLeaderboard && (
        <div className="modal-overlay">
          <div className="modal">
            <h2>🏆 Leaderboard ({difficulty})</h2>
            <table>
              <thead>
                <tr>
                  <th>Rank</th>
                  <th>Name</th>
                  <th>Time</th>
                </tr>
              </thead>
              <tbody>
                {leaderboard.length > 0 ? (
                  leaderboard.map((score, index) => (
                    <tr key={index}>
                      <td>{index + 1}</td>
                      <td>{score.username}</td>
                      <td>{formatTime(score.time_spent)}</td>
                    </tr>
                  ))
                ) : (
                  <tr><td colSpan="3">暂无数据</td></tr>
                )}
              </tbody>
            </table>
            <button onClick={() => setShowLeaderboard(false)} className="close-btn">Close</button>
          </div>
        </div>
      )}
    </div>
  );
}

export default App;