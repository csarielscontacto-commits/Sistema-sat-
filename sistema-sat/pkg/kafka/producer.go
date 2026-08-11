package kafka

import (
    "context"
    "encoding/json"
    "time"

    "github.com/segmentio/kafka-go"
    "github.com/sistema-sat/internal/models"
    "github.com/sistema-sat/pkg/logger"
)

type Producer struct {
    writer *kafka.Writer
    log    *logger.Logger
}

type ProducerConfig struct {
    Brokers   []string
    Topic     string
    BatchSize int
}

func NewProducer(cfg ProducerConfig, log *logger.Logger) *Producer {
    writer := &kafka.Writer{
        Addr:         kafka.TCP(cfg.Brokers...),
        Topic:        cfg.Topic,
        Balancer:     &kafka.LeastBytes{},
        BatchSize:    cfg.BatchSize,
        BatchTimeout: 100 * time.Millisecond,
        RequiredAcks: kafka.RequireAll,
        Compression:  kafka.Snappy,
        MaxAttempts:  5,
    }

    return &Producer{
        writer: writer,
        log:    log,
    }
}

func (p *Producer) PublishCFDI(ctx context.Context, cfdi *models.CFDI) error {
    data, err := json.Marshal(cfdi)
    if err != nil {
        return err
    }

    msg := kafka.Message{
        Key:   []byte(cfdi.UUID),
        Value: data,
        Time:  time.Now(),
        Headers: []kafka.Header{
            {Key: "event_type", Value: []byte("cfdi_ingested")},
            {Key: "rfc", Value: []byte(cfdi.Emisor.RFC)},
        },
    }

    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    err = p.writer.WriteMessages(ctx, msg)
    if err != nil {
        p.log.Errorw("Error publicando en Kafka",
            "topic", p.writer.Topic,
            "uuid", cfdi.UUID,
            "error", err,
        )
        return err
    }

    p.log.Debugw("CFDI publicado en Kafka",
        "topic", p.writer.Topic,
        "uuid", cfdi.UUID,
        "rfc", cfdi.Emisor.RFC,
    )

    return nil
}

func (p *Producer) PublishBatch(ctx context.Context, cfdis []*models.CFDI) error {
    messages := make([]kafka.Message, len(cfdis))

    for i, cfdi := range cfdis {
        data, err := json.Marshal(cfdi)
        if err != nil {
            return err
        }

        messages[i] = kafka.Message{
            Key:   []byte(cfdi.UUID),
            Value: data,
            Time:  time.Now(),
            Headers: []kafka.Header{
                {Key: "event_type", Value: []byte("cfdi_batch")},
            },
        }
    }

    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()

    err := p.writer.WriteMessages(ctx, messages...)
    if err != nil {
        p.log.Errorw("Error publicando batch en Kafka",
            "topic", p.writer.Topic,
            "count", len(cfdis),
            "error", err,
        )
        return err
    }

    return nil
}

func (p *Producer) Close() error {
    return p.writer.Close()
}

func (p *Producer) HealthCheck(ctx context.Context) error {
    testWriter := kafka.NewWriter(kafka.WriterConfig{
        Addr:  p.writer.Addr,
        Topic: "__health_check",
    })
    defer testWriter.Close()

    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    err := testWriter.WriteMessages(ctx, kafka.Message{
        Key:   []byte("health"),
        Value: []byte("ok"),
    })

    return err
}