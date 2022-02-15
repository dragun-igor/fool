import React from 'react'
import './styles.css'
import { GameCard } from '../GameCard'
import {Row, Col} from "antd";

export const Table = (props) => {
    const { pairs, onClick } = props
    let className = 'table'
    console.log(pairs)

    return (
        <div
            className={className}
            onClick={() => {onClick()}}
        >
            <h1 style={{
                marginBottom: 0
            }} align="center">Table</h1>
            <Row>
                {
                    pairs.map((pair) => {
                        return (
                            <Col>
                                <GameCard
                                    card={pair.first_card}
                                />
                                <GameCard
                                    card={pair.second_card}
                                />
                            </Col>
                        )
                    })
                }
            </Row>
        </div>
    )
}