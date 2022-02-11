import React, { useState } from 'react';
import { render } from 'react-dom';
import {DatePicker, message, Menu, Input, Button, Card, Row} from 'antd';
import 'antd/dist/antd.css';
import './index.css';
import ss from './4.png'

const App = () => {
    const [disabled, setDisabled] = React.useState(false);
    const toggle = () => {
        setDisabled(!disabled)
    }
    return (

            <Row>
                <div>
                <Menu
                    style={{ width: 300 }}
                    defaultSelectedKeys={['1']}
                    defaultOpenKeys={['sub1']}
                    theme="dark"
                    mode="inline"
                >
                    <Menu.SubMenu key="sub1" title={
                        <span>
                            <span>Menu</span>
                        </span>}
                    >
                        <Menu.ItemGroup key="g1" title="Select Game">
                            <Menu.Item key="1">1</Menu.Item>
                            <Menu.Item key="2">2</Menu.Item>
                            <Menu.Item key="3">3</Menu.Item>
                            <Menu.Item key="4">4</Menu.Item>
                            <Menu.Item key="5">5</Menu.Item>
                            <Menu.Item key="6">6</Menu.Item>
                            <Menu.Item key="7">7</Menu.Item>
                            <Menu.Item key="8">8</Menu.Item>
                            <Menu.Item key="9">9</Menu.Item>
                        </Menu.ItemGroup>
                    </Menu.SubMenu>
                </Menu>
                </div>
                <div>
                    <h1 align='center'>Hand</h1>
                    <Row>
                        <Card
                            hoverable
                            cover={
                                <img
                                    alt="example"
                                    src={ss}
                                />
                            }
                            style={{ width: 120, height: 300, borderColor: 'rgba(255, 215, 0)', borderWidth: 2 }}
                        >
                            <Card.Meta
                                title="Six Spades"
                                description="Trump Suit" />
                        </Card>
                        <Card
                            hoverable
                            cover={
                                <img
                                    alt="six spades"
                                    src={ss}
                                />
                            }
                            style={{
                                width: 120,
                                height:180,
                                borderColor: 'rgba(255, 215, 0)',
                                borderWidth: 2,
                                borderRadius: 10
                            }}
                        />
                        <Card
                            hoverable
                            title="Card"
                            style={{ width: 150 }}
                        >
                            <Card.Meta
                                title="Six Spades"
                                description="Not Trump Suit" />
                        </Card>
                    </Row>
                    <Input.Group compact style={{ margin: 10 }}>
                        <Input style={{ width: 222 }} placeholder='input your name' disabled={disabled} />
                        <Button type='primary' shape='round' onClick={toggle}>Submit</Button>
                    </Input.Group>
                </div>




                <div>
                    <Card
                        title="Hand"
                        >
                        <Row>
                            <Card
                                hoverable
                                cover={
                                    <img
                                        alt="six spades"
                                        src={ss}
                                    />
                                }
                                style={{
                                    width: 120,
                                    height:180,
                                    borderColor: 'rgba(255, 215, 0)',
                                    borderWidth: 2,
                                    borderRadius: 10
                                }}
                            />
                            <Button type='ghost' style={{
                                width: 120,
                                height: 180,
                            }}> </Button>
                            <Card
                                hoverable
                                cover={
                                    <img
                                        alt="six spades"
                                        src={ss}
                                    />
                                }
                                style={{
                                    width: 120,
                                    height:180,
                                    borderColor: 'rgba(255, 215, 0)',
                                    borderWidth: 2,
                                    borderRadius: 10
                                }}
                            />
                            <Card
                                hoverable
                                cover={
                                    <img
                                        alt="six spades"
                                        src={ss}
                                    />
                                }
                                style={{
                                    width: 120,
                                    height:180,
                                    borderColor: 'rgba(255, 215, 0)',
                                    borderWidth: 2,
                                    borderRadius: 10
                                }}
                            />
                            <Card
                                hoverable
                                cover={
                                    <img
                                        alt="six spades"
                                        src={ss}
                                    />
                                }
                                style={{
                                    width: 120,
                                    height:180,
                                    borderRadius: 10
                                }}
                            />
                            <Card
                                hoverable
                                cover={
                                    <img
                                        alt="six spades"
                                        src={ss}
                                    />
                                }
                                style={{
                                    width: 120,
                                    height:180,
                                    borderColor: 'rgba(255, 215, 0)',
                                    borderWidth: 2,
                                    borderRadius: 10
                                }}
                            />
                            <Card
                                hoverable
                                cover={
                                    <img
                                        alt="six spades"
                                        src={ss}
                                    />
                                }
                                style={{
                                    width: 120,
                                    height: 180,
                                    borderRadius: 10
                                }}
                            />
                        </Row>
                    </Card>
                </div>




            </Row>

    );
};
render(<App />, document.getElementById('root'));
