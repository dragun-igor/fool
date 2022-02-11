import logo from './logo.svg';
// import React, { useState } from 'react';
// import { render } from 'react-dom';
// import { DatePicker, message } from 'antd';
// import 'antd/dist/antd.css';
// import './index.css';

function App() {
  return (
    <div>
      <Menu
        style={{ width: 256 }}
        defaultSelectedKeys={['1']}
        defaultOpenKeys={['sub1']}
        mode="inline"
      >
        <SubMenu key="sub1" title={
          <span>
            <span>Navigation One</span>
          </span>}
        >
          <MenuItemGroup key="g1" title="Item 1">
            <Menu.Item key="1">Option 1</Menu.Item>
            <Menu.Item key="2">Option 2</Menu.Item>
          </MenuItemGroup>
        </SubMenu>
      </Menu>
    </div>
  );
}

export default App;

// const App = () => {
//   const [date, setDate] = useState(null);
//   const handleChange = value => {
//     message.info(`Selected Date: ${value ? value.format('YYYY-MM-DD') : 'None'}`);
//     setDate(value);
//   };
//   return (
//       <div style={{ width: 400, margin: '100px auto' }}>
//         <DatePicker onChange={handleChange} />
//         <div style={{ marginTop: 16 }}>
//           Selected Date: {date ? date.format('YYYY-MM-DD') : 'None'}
//         </div>
//       </div>
//   );
// };

// render(<App />, document.getElementById('root'));
