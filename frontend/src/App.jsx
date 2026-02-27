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

  // --- 新增状态 ---
  const [isEligible, setIsEligible] = useState(true); // 是否有资格提交成绩
  const [conflicts, setConflicts] = useState(new Set()); // 存储错误的格子坐标 "row-col"

  const timerRef = useRef(null);

  const formatTime = (seconds) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
  };

  const stopTimer = () => {
    if (timerRef.current) {
      clearInterval(timerRef.current);
      timerRef.current = null;
    }
  };

  const fetchNewGame = async (level) => {
    stopTimer();
    setLoading(true);
    setIsGameActive(false);
    setIsWon(false);
    setTimer(0);
    setShowLeaderboard(false);
    
    // 重置状态
    setIsEligible(true);
    setConflicts(new Set());

    try {
      const res = await axios.get(`${API_URL}/api/new-game?level=${level}`);
      
      stopTimer();

      setBoard(res.data.puzzle);
      setInitialBoard(JSON.parse(JSON.stringify(res.data.puzzle)));
      setSolution(res.data.solution); 
      
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
    // 允许输入空值（删除）或 1-9
    if (value === '' || (num >= 1 && num <= 9)) {
      const newBoard = JSON.parse(JSON.stringify(board));
      newBoard[rowIndex][colIndex] = value === '' ? 0 : num;
      setBoard(newBoard);
      
      // 用户修改了格子，移除该格子的错误高亮
      const key = `${rowIndex}-${colIndex}`;
      if (conflicts.has(key)) {
        const newConflicts = new Set(conflicts);
        newConflicts.delete(key);
        setConflicts(newConflicts);
      }

      const hasEmpty = newBoard.some(row => row.includes(0));
      if (!hasEmpty) {
        checkWin(newBoard);
      }
    }
  };

  // --- 新功能：AI 提示 ---
  const handleHint = () => {
    if (!isGameActive) return;
    
    // 标记成绩无效
    setIsEligible(false);

    // 找到所有空格子
    const emptySpots = [];
    board.forEach((row, r) => {
      row.forEach((val, c) => {
        if (val === 0) emptySpots.push({ r, c });
      });
    });

    if (emptySpots.length === 0) return;

    // 随机选一个空格
    const randomSpot = emptySpots[Math.floor(Math.random() * emptySpots.length)];
    const { r, c } = randomSpot;

    // 填入正确答案
    const newBoard = JSON.parse(JSON.stringify(board));
    newBoard[r][c] = solution[r][c];
    setBoard(newBoard);

    // 检查是否获胜
    const hasEmpty = newBoard.some(row => row.includes(0));
    if (!hasEmpty) checkWin(newBoard);
  };

  // --- 新功能：查看答案 ---
  const handleSolve = () => {
    if (!isGameActive) return;
    if (!window.confirm("确定要查看答案吗？这将无法提交成绩。")) return;

    setIsEligible(false);
    setBoard(JSON.parse(JSON.stringify(solution))); // 直接填满
    stopTimer();
    setIsGameActive(false);
    // 注意：这里我们不触发 setIsWon，因为这是放弃比赛
  };

  // --- 新功能：检查冲突 ---
  const handleCheck = () => {
    if (!isGameActive) return;

    const newConflicts = new Set();
    board.forEach((row, r) => {
      row.forEach((val, c) => {
        // 如果格子填了数字，且不等于答案，就是错误的
        if (val !== 0 && val !== solution[r][c]) {
          newConflicts.add(`${r}-${c}`);
        }
      });
    });

    setConflicts(newConflicts);
    
    // 3秒后自动清除高亮（可选，提升体验）
    if (newConflicts.size > 0) {
      setTimeout(() => setConflicts(new Set()), 3000);
    }
  };

  const submitScore = async () => {
    if (!username) return alert("请输入名字！");
    try {
      await axios.post(`${API_URL}/api/submit-score`, {
        username,
        difficulty,
        time_spent: timer 
      });
      alert("分数提交成功！");
      setIsWon(false);
      fetchLeaderboard(difficulty);
      setShowLeaderboard(true);
    } catch (error) {
      console.error("Submit Error:", error);
      alert("提交失败");
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
    return () => stopTimer();
  }, []);

  return (
    <div className="container">
      <h1>Sudoku Go</h1>
      
      <div className="header-info">
        <div className="timer">⏱️ {formatTime(timer)}</div>
        <button className="btn-secondary" onClick={() => {
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

      {/* 新增工具栏 */}
      <div className="tools">
        <button className="btn-tool" onClick={handleCheck}>🔍 Check</button>
        <button className="btn-tool" onClick={handleHint}>💡 Hint</button>
        <button className="btn-tool btn-danger" onClick={handleSolve}>👁️ Solve</button>
      </div>
      {!isEligible && <div className="warning-text">⚠️ 辅助功能已使用，本局成绩无效</div>}

      <div className="board">
        {board.map((row, rowIndex) => (
          <div key={rowIndex} className="row">
            {row.map((cell, colIndex) => {
              const isInitial = initialBoard[rowIndex][colIndex] !== 0;
              const isConflict = conflicts.has(`${rowIndex}-${colIndex}`);
              return (
                <input
                  key={`${rowIndex}-${colIndex}`}
                  type="text"
                  maxLength="1"
                  // 动态添加 conflict 类
                  className={`cell ${isInitial ? 'initial' : ''} ${isConflict ? 'conflict' : ''}`}
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
            <p>Time: {formatTime(timer)}</p>
            
            {/* 只有 isEligible 为 true 时才允许提交 */}
            {isEligible ? (
              <>
                <input 
                  type="text" 
                  placeholder="Enter your name" 
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                />
                <button onClick={submitScore}>Submit Score</button>
              </>
            ) : (
              <p className="error-msg">辅助功能已使用，无法提交成绩。</p>
            )}
            
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