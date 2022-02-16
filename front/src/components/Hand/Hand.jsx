import React from 'react'
import './styles.css'
import { GameCard } from '../GameCard'
import { Row, Col } from 'antd'

export const Hand = (props) => {
    const {
        cards,
        onClick
    } = props
    let className = 'hand'

    return (
        <div className={className}>
            <h1 style={{
                marginBottom: 0
            }} align="center">Your Hand</h1>
            <Row>
                {
                    cards.map((card) => {
                        return (
                            <GameCard
                                card={card}
                                onClick={onClick}
                            />
                        )
                    })
                }
            </Row>
        </div>
    )
}