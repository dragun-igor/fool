import React, {useEffect, useState} from 'react'
import { Hand } from './components/Hand'
import { Table } from './components/Table'
import { Input, Button } from 'antd'

const socket = new WebSocket('ws://127.0.0.1:8080/ws')

function App() {
    const [disabled, setDisabled] = useState(false)
    const toggle = () => {
        setDisabled(!disabled)
    }
    const [inputUsernameField, setInputUsernameField] = useState('')
    const UsernameField = Input

    const [state, setState] = useState({
        connected_users: [],
        table: [],
        hand: []
    })

    const changeTable = (pairs) => {
        setState({
            ...state,
            table: pairs
        })
    }

    const changeHand = (cards) => {
        setState({
            ...state,
            hand: cards
        })
    }

    const changeUserList = (users) => {
        setState({
            ...state,
            connected_users: users
        })
    }

    const handleSetUsername = () => {
        let jsonData = {}
        jsonData.action = "add_user"
        jsonData.username = inputUsernameField
        socket.send(JSON.stringify(jsonData))
    }

    const handleStartGame = () => {
        let jsonData = {}
        jsonData.action = "start_game"
        socket.send(JSON.stringify(jsonData))
        toggle()
    }

    const onClickHandCard = (id) => {
        let jsonData = {}
        jsonData.action = "select_card"
        jsonData.id = id
        socket.send(JSON.stringify(jsonData))
    }

    const onClickTable = () => {
        let jsonData = {}
        jsonData.action = "put_on_table"
        socket.send(JSON.stringify(jsonData))
    }

    useEffect(() => {
        socket.onmessage = msg => {
            let data = JSON.parse(msg.data)

            switch (data.action) {
                case "list_users":
                    changeUserList(data.connected_users)
                    break
                case "hand":
                    changeHand(data.hand)
                    break
                case "table":
                    changeTable(data.table.pairs)
                    break
            }
        }
    })

  return (
      <div>
          <h1>{state.connected_users}</h1>
          <Input.Group compact style={{
              margin: 10
          }}>
              <UsernameField
                  style={{
                      width: 222
                  }}
                  placeholder="input your name"
                  disabled={disabled}
                  onChange={(e) => setInputUsernameField(e.target.value)}
                  onPressEnter={handleSetUsername}
              />

              <Button
                  type='primary'
                  disabled={disabled}
                  onClick={handleSetUsername}
              >
                  Submit
              </Button>
              <Button
                  type='default'
                  disabled={disabled}
                  onClick={handleStartGame}
              >
                  Start
              </Button>
          </Input.Group>
          <Hand
              cards={state.hand}
              onClick={onClickHandCard}
          />
          <Table
              pairs={state.table}
              onClick={onClickTable}
          />
      </div>
  )
}

export default App;


