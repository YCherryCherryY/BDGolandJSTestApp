
import './App.css';
import React from 'react'
import { Layout } from 'antd';
import StatusTable from './components/StatusTable';

const {Header, Content} = Layout;

function App() {
  return (
    <Layout>
      <Header>
        <h1 style={{ color: 'white' }}>Container Status Monitor</h1>
      </Header>
      <Content style={{ padding: '20px' }}>
        <StatusTable />
      </Content>
    </Layout>
  );
}

export default App;
