import React, { useEffect, useState } from 'react';
import { BrowserRouter as Router, Routes, Route, useNavigate, useParams } from 'react-router-dom';

const API_BASE = "https://order-coffee-production.up.railway.app";

function CoffeeCard({ coffee, onClick, quantity, setQuantity }) {
  return (
    <div style={{ border: '1px solid #ccc', borderRadius: 8, padding: 16, margin: 8, width: 220, textAlign: 'center', boxShadow: '0 2px 8px #eee', cursor: 'pointer' }} onClick={onClick}>
      <img src={`/images/${coffee.image}`} alt={coffee.name} style={{ width: '100%', height: 120, objectFit: 'cover', borderRadius: 8 }} />
      <h2 style={{ margin: '12px 0 8px 0' }}>{coffee.name}</h2>
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', marginBottom: 8 }}>
        <button type="button" onClick={e => { e.stopPropagation(); setQuantity(Math.max(1, quantity - 1)); }} style={{ width: 32, height: 32 }}>-</button>
        <span style={{ margin: '0 12px' }}>{quantity}</span>
        <button type="button" onClick={e => { e.stopPropagation(); setQuantity(quantity + 1); }} style={{ width: 32, height: 32 }}>+</button>
      </div>
    </div>
  );
}

function CoffeeList() {
  const [coffees, setCoffees] = useState([]);
  const [quantities, setQuantities] = useState({});
  const navigate = useNavigate();

  useEffect(() => {
    fetch(`${API_BASE}/coffees`)
      .then(res => res.json())
      .then(data => {
        setCoffees(data);
        // Set default quantity to 1 for each coffee
        const q = {};
        data.forEach(c => { q[c.id] = 1; });
        setQuantities(q);
      });
  }, []);

  return (
    <div style={{ maxWidth: 900, margin: '2rem auto', fontFamily: 'sans-serif' }}>
      <h1>Order Coffee</h1>
      <div style={{ display: 'flex', flexWrap: 'wrap', justifyContent: 'center' }}>
        {coffees.map(coffee => (
          <CoffeeCard
            key={coffee.id}
            coffee={coffee}
            quantity={quantities[coffee.id] || 1}
            setQuantity={q => setQuantities({ ...quantities, [coffee.id]: q })}
            onClick={() => navigate(`/order/${coffee.id}?quantity=${quantities[coffee.id] || 1}`)}
          />
        ))}
      </div>
    </div>
  );
}

function OrderPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [coffee, setCoffee] = useState(null);
  const [quantity, setQuantity] = useState(1);
  const [name, setName] = useState('');
  const [notes, setNotes] = useState('');
  const [message, setMessage] = useState('');

  useEffect(() => {
    fetch(`${API_BASE}/coffees`)
      .then(res => res.json())
      .then(data => {
        const c = data.find(c => String(c.id) === id);
        setCoffee(c);
      });
    // Get quantity from query param
    const params = new URLSearchParams(window.location.search);
    setQuantity(parseInt(params.get('quantity')) || 1);
  }, [id]);

  const handleSubmit = (e) => {
    e.preventDefault();
    fetch(`${API_BASE}/order`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        name,
        coffeeId: parseInt(id),
        notes: notes + ` (Quantity: ${quantity})`
      })
    })
      .then(res => res.json())
      .then(data => {
        setMessage('Order submitted! Redirecting...');
        setTimeout(() => navigate('/'), 1500);
      });
  };

  if (!coffee) return <div>Loading...</div>;

  return (
    <div style={{ maxWidth: 400, margin: '2rem auto', fontFamily: 'sans-serif' }}>
      <img src={`/images/${coffee.image}`} alt={coffee.name} style={{ width: '100%', height: 180, objectFit: 'cover', borderRadius: 8 }} />
      <h2>{coffee.name}</h2>
      <p>Quantity: {quantity}</p>
      <form onSubmit={handleSubmit}>
        <input
          placeholder="Your name"
          value={name}
          onChange={e => setName(e.target.value)}
          required
          style={{ width: '100%', marginBottom: 8, padding: 8 }}
        />
        <input
          placeholder="Notes"
          value={notes}
          onChange={e => setNotes(e.target.value)}
          style={{ width: '100%', marginBottom: 8, padding: 8 }}
        />
        <button type="submit" style={{ width: '100%', padding: 8 }}>Send Order</button>
      </form>
      {message && <p style={{ color: 'green' }}>{message}</p>}
    </div>
  );
}

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<CoffeeList />} />
        <Route path="/order/:id" element={<OrderPage />} />
      </Routes>
    </Router>
  );
}

export default App;
